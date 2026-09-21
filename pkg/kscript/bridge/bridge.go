// Package bridge runs legacy Kite script tasks through the standalone
// github.com/gookit/taskrun library.
//
// The legacy runner in pkg/kscript stays in place as the rollback point. This
// package only converts what the legacy configuration already loaded and runs a
// single task or script file with the new engine; command aliases, extensions,
// plugins and system command fallback remain owned by Kite.
package bridge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/slog"
	"github.com/gookit/taskrun"
	"github.com/gookit/taskrun/formats"

	"github.com/inhere/kite-go/pkg/kscript"
)

// BaseVarNames lists the runtime variable names Kite injects in addition to the
// caller supplied variables. They mirror the legacy render context.
var BaseVarNames = []string{
	"gvs", "paths", "kite", "time", "workdir", "dirname", "cur_dir", "@", "*",
}

// Bridge converts the legacy runner configuration and runs tasks with the new
// engine.
type Bridge struct {
	legacy   *kscript.Runner
	baseDir  string
	extras   []string
	envNames []string
	parseEnv bool
	handlers map[string]taskrun.Handler
	io       taskrun.IO
	fallback bool

	mu    sync.Mutex
	cache map[string]*entry
}

type entry struct {
	runner   *taskrun.Runner
	warnings []string
}

// Option configures a Bridge.
type Option func(*Bridge)

// WithBaseDir sets the absolute directory used by the converted definition.
// It defaults to the current working directory, matching the legacy behavior of
// resolving relative paths against the process directory.
func WithBaseDir(dir string) Option {
	return func(b *Bridge) {
		if dir != "" {
			b.baseDir = dir
		}
	}
}

// WithVarNames adds runtime variable names beyond BaseVarNames, so legacy
// $name and ${name} references to them are rewritten.
func WithVarNames(names ...string) Option {
	return func(b *Bridge) {
		b.extras = append(b.extras, names...)
	}
}

// WithHandler registers a host action for converted definitions.
func WithHandler(name string, handler taskrun.Handler) Option {
	return func(b *Bridge) {
		if b.handlers == nil {
			b.handlers = map[string]taskrun.Handler{}
		}
		b.handlers[name] = handler
	}
}

// WithIO sets the streams used for runs. It defaults to the process streams.
func WithIO(io taskrun.IO) Option {
	return func(b *Bridge) { b.io = io }
}

// WithLegacyFallback makes TryRun fall back to the legacy runner when the
// conversion fails, so a partially converted configuration never regresses. It
// logs the conversion error at warning level.
func WithLegacyFallback(enabled bool) Option {
	return func(b *Bridge) { b.fallback = enabled }
}

// WithParseEnv enables the legacy ParseEnv behavior, where a bare $NAME that
// matches an environment variable is replaced by its value at render time.
// When disabled, such references are left for the shell to expand.
func WithParseEnv(enabled bool) Option {
	return func(b *Bridge) { b.parseEnv = enabled }
}

// New creates a bridge over an initialized or not yet initialized legacy
// runner.
func New(legacy *kscript.Runner, opts ...Option) *Bridge {
	b := &Bridge{
		legacy:   legacy,
		cache:    map[string]*entry{},
		parseEnv: legacy != nil && legacy.ParseEnv,
	}
	for _, opt := range opts {
		opt(b)
	}
	if b.baseDir == "" {
		b.baseDir = sysutil.Workdir()
	}
	if b.io.Stdout == nil {
		b.io.Stdout = os.Stdout
	}
	if b.io.Stderr == nil {
		b.io.Stderr = os.Stderr
	}
	return b
}

// Definition converts the loaded legacy configuration. The shell argument is
// the legacy per-run shell wrapper (RunCtx.Type); an empty value keeps the
// definition time conversion.
func (b *Bridge) Definition(shell string, varNames, envNames []string) (taskrun.Definition, []string, error) {
	if err := b.legacy.InitLoad(); err != nil {
		return taskrun.Definition{}, nil, err
	}
	extToBin := b.legacy.ExtToBin()
	files := map[string]formats.LegacyScriptFile{}
	for name, path := range b.legacy.ScriptFileMap() {
		ext := filepath.Ext(name)
		if name == "" || path == "" {
			continue
		}
		files[name] = formats.LegacyScriptFile{Path: path, Ext: ext, Bin: extToBin[ext]}
	}
	result, err := formats.LegacyDefinition(formats.LegacyOptions{
		Scripts:      b.legacy.Scripts,
		Settings:     b.legacy.SettingsData(),
		BaseDir:      b.baseDir,
		DefaultShell: shell,
		RuntimeVars:  varNames,
		EnvNames:     envNames,
		Files:        files,
	})
	if err != nil {
		return taskrun.Definition{}, result.Warnings, err
	}
	return result.Definition, result.Warnings, nil
}

// runnerFor returns the converted runner for a shell and runtime variable set.
// Definitions are cached because conversion and validation are pure.
func (b *Bridge) runnerFor(shell string, varNames, envNames []string) (*entry, error) {
	key := shell + "\x00" + strings.Join(varNames, ",") + "\x00" + strings.Join(envNames, ",")
	b.mu.Lock()
	defer b.mu.Unlock()
	if cached, ok := b.cache[key]; ok {
		return cached, nil
	}
	def, warnings, err := b.Definition(shell, varNames, envNames)
	if err != nil {
		return nil, err
	}
	opts := make([]taskrun.Option, 0, len(b.handlers))
	for name, handler := range b.handlers {
		opts = append(opts, taskrun.WithHandler(name, handler))
	}
	runner, err := taskrun.New(def, opts...)
	if err != nil {
		return nil, err
	}
	cached := &entry{runner: runner, warnings: warnings}
	b.cache[key] = cached
	return cached, nil
}

// TryRun runs a script task or script file by name. It reports found=false only
// when the name matches neither, mirroring the legacy TryRun contract so the
// caller can continue with system command fallback.
func (b *Bridge) TryRun(ctx context.Context, name string, args []string, lctx *kscript.RunCtx) (bool, error) {
	if err := b.legacy.InitLoad(); err != nil {
		return true, err
	}
	lctx = kscript.EnsureCtx(lctx).WithNameArgs(name, args)
	if strings.ContainsRune(name, ' ') {
		// Keyword search is Kite CLI sugar and stays with the legacy runner.
		return b.legacyTryRun(name, args, lctx)
	}
	workdir := b.workdir(lctx)
	vars := b.runtimeVars(lctx, args, workdir)
	names := b.runtimeVarNames(vars)
	envNames := b.envNameList(lctx)
	cached, err := b.runnerFor(lctx.Type, names, envNames)
	if err != nil {
		if b.fallback {
			slog.Warnf("kscript bridge: conversion failed, falling back to the legacy runner: %v", err)
			return b.legacyTryRun(name, args, lctx)
		}
		return true, err
	}
	for _, warning := range cached.warnings {
		slog.Debugf("kscript bridge: %s", warning)
	}
	if _, err := cached.runner.Lookup(name); err == nil {
		if !lctx.Silent {
			slog.Infof("Run Script Task %q Args: %v (engine: kscript)", name, args)
		}
		return true, b.runTask(ctx, cached.runner, name, args, lctx, vars, workdir)
	}
	// A discovered script file is runnable by its file name.
	if file, ok := b.scriptFile(name); ok {
		if !lctx.Silent {
			slog.Infof("Run Script File %q Args: %v (engine: kscript)", name, args)
		}
		return true, b.runFile(ctx, file, name, args, lctx, vars, workdir)
	}
	return false, nil
}

// Run is TryRun with the legacy error message when nothing matched.
func (b *Bridge) Run(ctx context.Context, name string, args []string, lctx *kscript.RunCtx) error {
	found, err := b.TryRun(ctx, name, args, lctx)
	if !found {
		return fmt.Errorf("script file %q is not exists", name)
	}
	return err
}

func (b *Bridge) runTask(ctx context.Context, runner *taskrun.Runner, name string, args []string, lctx *kscript.RunCtx, vars map[string]any, workdir string) error {
	request := taskrun.Request{
		Task:   name,
		Args:   args,
		Vars:   vars,
		Env:    lctx.Env,
		Dir:    workdir,
		DryRun: lctx.DryRun,
		IO:     b.io,
	}
	if lctx.DryRun {
		plan, err := runner.Inspect(ctx, request)
		if err != nil {
			return err
		}
		for _, action := range plan.Actions {
			fmt.Fprintf(b.io.Stdout, "[dry-run] %s %s %s %v\n", action.Task, action.Kind, action.Program, action.Args)
		}
		for _, deferred := range plan.Deferred {
			fmt.Fprintf(b.io.Stdout, "[dry-run] deferred: %s\n", deferred)
		}
		return nil
	}
	result, err := runner.Run(ctx, request)
	if err != nil {
		return err
	}
	if result != nil {
		slog.Debugf("kscript: task %s status=%s steps=%d", name, result.Status, len(result.Steps))
	}
	return nil
}

// runFile builds a one-off definition for a discovered script file, because the
// arguments of a file action are known only at run time.
func (b *Bridge) runFile(ctx context.Context, file formats.LegacyScriptFile, name string, args []string, lctx *kscript.RunCtx, vars map[string]any, workdir string) error {
	program, prefix, err := formats.SplitCommandLine(file.Bin)
	if err != nil {
		return err
	}
	if program == "" {
		return fmt.Errorf("script file %q has no interpreter", name)
	}
	fileKey := "file"
	def := taskrun.Definition{
		Version: 1,
		BaseDir: b.baseDir,
		Files: map[string]taskrun.ScriptFile{
			fileKey: {
				Name:        fileKey,
				Path:        file.Path,
				Interpreter: taskrun.Interpreter{Program: program, PrefixArgs: prefix},
			},
		},
		Tasks: map[string]taskrun.Task{
			name: {
				Name:  name,
				Steps: []taskrun.Step{{Name: name, File: &taskrun.FileSpec{Name: fileKey, Args: args}}},
			},
		},
	}
	runner, err := taskrun.New(def)
	if err != nil {
		return err
	}
	return b.runTask(ctx, runner, name, args, lctx, vars, workdir)
}

// scriptFile resolves a discovered script file the way the legacy runner does:
// an explicit extension wins, otherwise the allowed extensions are probed.
func (b *Bridge) scriptFile(name string) (formats.LegacyScriptFile, bool) {
	extToBin := b.legacy.ExtToBin()
	files := b.legacy.ScriptFileMap()
	build := func(key string) (formats.LegacyScriptFile, bool) {
		path, ok := files[key]
		if !ok {
			return formats.LegacyScriptFile{}, false
		}
		ext := filepath.Ext(key)
		return formats.LegacyScriptFile{Path: path, Ext: ext, Bin: extToBin[ext]}, true
	}
	if filepath.Ext(name) != "" {
		return build(name)
	}
	for _, ext := range b.legacy.AllowedExt {
		if file, ok := build(name + ext); ok {
			return file, true
		}
	}
	return formats.LegacyScriptFile{}, false
}

func (b *Bridge) legacyTryRun(name string, args []string, lctx *kscript.RunCtx) (bool, error) {
	return b.legacy.TryRun(name, args, lctx)
}

func (b *Bridge) workdir(lctx *kscript.RunCtx) string {
	if lctx != nil && lctx.Workdir != "" {
		return lctx.Workdir
	}
	return sysutil.Workdir()
}

// runtimeVars mirrors the legacy render context: caller variables, the argument
// shorthands ($@, $*, $1..$N), run metadata and the AppendVarsFn extensions.
func (b *Bridge) runtimeVars(lctx *kscript.RunCtx, args []string, workdir string) map[string]any {
	vars := map[string]any{}
	for key, value := range lctx.Vars {
		vars[key] = value
	}
	argStr := cliutil.BuildLine("", args)
	vars["@"] = argStr
	vars["*"] = strutil.Quote(argStr)
	for i, value := range args {
		vars[strconv.Itoa(i+1)] = value
	}
	now := time.Now()
	vars["time"] = map[string]any{
		"unix_sec":    now.Unix(),
		"datetime":    now.Format("2006-01-02 15:04:05"),
		"date_Ymd_hm": now.Format("2006-01-02_15:04"),
		"date_ymd_hm": now.Format("06-01-02_15:04"),
		"date_ymd":    now.Format("2006-01-02"),
		"date_hms":    now.Format("15:04:05"),
	}
	vars["cur_dir"] = sysutil.Workdir()
	vars["workdir"] = workdir
	vars["dirname"] = fsutil.Name(workdir)
	if lctx.AppendVarsFn != nil {
		vars = lctx.AppendVarsFn(vars)
	}
	return vars
}

// runtimeVarNames lists every name a template may reference, so the converter
// rewrites the matching legacy references.
func (b *Bridge) runtimeVarNames(vars map[string]any) []string {
	seen := map[string]bool{}
	for _, name := range BaseVarNames {
		seen[name] = true
	}
	for _, name := range b.extras {
		if name != "" {
			seen[name] = true
		}
	}
	for name := range vars {
		seen[name] = true
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// envNameList lists the environment names the legacy ParseEnv option could
// substitute. It is empty unless ParseEnv is enabled, in which case the process
// environment, the run context environment and any configured extras are
// included.
func (b *Bridge) envNameList(lctx *kscript.RunCtx) []string {
	if !b.parseEnv {
		return nil
	}
	seen := map[string]bool{}
	for _, item := range os.Environ() {
		if index := strings.IndexByte(item, '='); index > 0 {
			seen[item[:index]] = true
		}
	}
	if lctx != nil {
		for name := range lctx.Env {
			seen[name] = true
		}
	}
	for _, name := range b.envNames {
		seen[name] = true
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// WithEnvNames adds environment names for ParseEnv substitution.
func WithEnvNames(names ...string) Option {
	return func(b *Bridge) {
		b.envNames = append(b.envNames, names...)
	}
}
