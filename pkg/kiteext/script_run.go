package kiteext

import (
	"os"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/pkg/cmdutil"
)

// RunCtx struct
type RunCtx struct {
	// Name for script run
	Name string
	Type string

	// DryRun script
	DryRun bool
	// Workdir for run script
	Workdir string
	// Env setting for run
	Env map[string]string

	// BeforeFn hook
	BeforeFn func(si *ScriptItem)
}

// EnsureCtx to
func EnsureCtx(ctx *RunCtx) *RunCtx {
	if ctx == nil {
		return &RunCtx{}
	}
	return ctx
}

func (c *RunCtx) WithName(name string) *RunCtx {
	c.Name = name
	return c
}

// ScriptItem struct
type ScriptItem struct {
	// Type wrap for run script. allow: sh, bash, zsh
	Type string

	// Workdir for run script
	Workdir string

	// Name for script
	Name string
	// Desc message
	Desc string
	// Env setting for run
	Env map[string]string
	// Args script args definition.
	Args, ArgNames []string
	// Cmds commands define in ScriptRunner.DefineFiles
	Cmds []string

	// File script file path in ScriptRunner.ScriptDirs
	File string
	Bin  string
	Ext  string // eg: .go
}

// args type: string, strings
func (si *ScriptItem) loadArgsDefine(args any) error {
	if args == nil {
		return nil
	}

	switch typVal := args.(type) {
	case string: // desc
		si.Args = []string{typVal}
	case []string: // desc list
		si.Args = typVal
	case []any: // desc list
		si.Args = arrutil.SliceToStrings(typVal)
	// case map[string]string: // name with desc TODO map cannot be ordered
	// 	si.Args = typVal
	default:
		return errorx.Rawf("invalid args config for %q", si.Name)
	}
	return nil
}

// ScriptRunner struct
type ScriptRunner struct {
	// Scripts config and loaded from DefineFiles.
	//
	// format: {name: info, name2: info2, ...}
	Scripts map[string]any `json:"scripts"`
	// DefineFiles scripts define files, will read and add to Scripts
	DefineFiles []string `json:"define_files"`
	// WrapShell for run each script.
	//
	// value like: bash, sh, zsh or empty for direct run command
	WrapShell string `json:"wrap_shell"`
	// ParseEnv var on script command
	ParseEnv bool `json:"parse_env"`

	// ScriptDirs script file dirs, allow multi
	ScriptDirs []string `json:"script_dirs"`
	// AllowedExt allowed script file extensions. eg: .go, .sh
	AllowedExt []string `json:"allowed_ext"`
	// FindBinByExt on run a script file
	FindBinByExt bool `json:"find_bin_by_ext"`
	// ExtToBinMap settings
	ExtToBinMap map[string]string `json:"ext_to_bin_map"`
	// PathResolver handler
	PathResolver func(path string) string

	// loaded from ScriptDirs.
	//
	// format: {filename: filepath, ...}
	scriptFiles map[string]string
	// mark script loaded
	fileLoaded   bool
	defineLoaded bool
}

// AllowTypes for run script
var AllowTypes = []string{"sh", "zsh", "bash"}

// AllowExt list
var AllowExt = []string{".sh", ".zsh", ".bash", ".php", ".go", ".gop", ".kts", ".java", ".gry", ".groovy"}

// ExtToBinMap data
//
// eg:
//
//	'#!/usr/bin/env bash'
//	'#!/usr/bin/env -S go run'
var ExtToBinMap = map[string]string{
	".sh":     "sh",
	".zsh":    "zsh",
	".bash":   "bash",
	".php":    "php",
	".gry":    "groovy",
	".groovy": "groovy",
	".go":     "go run",
}

// NewScriptRunner instance
func NewScriptRunner(fns ...func(sr *ScriptRunner)) *ScriptRunner {
	sr := &ScriptRunner{
		ParseEnv:     true,
		AllowedExt:   AllowExt,
		ExtToBinMap:  ExtToBinMap,
		PathResolver: sysutil.ExpandPath,
		scriptFiles:  map[string]string{},
	}

	for _, fn := range fns {
		fn(sr)
	}
	return sr
}

/*
-----------
--------------------------------- Init load ---------------------------------
-----------
*/

// InitLoad define scripts and script files.
func (sr *ScriptRunner) InitLoad() error {
	if err := sr.LoadDefineScripts(); err != nil {
		return err
	}
	return sr.LoadScriptFiles()
}

// LoadDefineScripts from DefineFiles
func (sr *ScriptRunner) LoadDefineScripts() error {
	if sr.defineLoaded {
		return nil
	}

	sr.defineLoaded = true
	loader := config.New("loader")

	for _, fpath := range sr.DefineFiles {
		fpath = sr.PathResolver(fpath)

		err := loader.LoadFiles(fpath)
		if err != nil {
			return err
		}

		sr.Scripts = maputil.SimpleMerge(loader.Data(), sr.Scripts)
		loader.ClearData()
	}

	return nil
}

// LoadScriptFiles from the ScriptDirs
func (sr *ScriptRunner) LoadScriptFiles() error {
	if sr.fileLoaded {
		return nil
	}
	sr.fileLoaded = true

	for _, dirPath := range sr.ScriptDirs {
		dirPath = sr.PathResolver(dirPath)
		des, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, ent := range des {
			fName := ent.Name()
			if !ent.IsDir() {
				// add
				sr.scriptFiles[fName] = dirPath + "/" + fName
			}
		}
	}

	return nil
}

/*
-----------
--------------------------------- Run script ---------------------------------
-----------
*/

// Run script or script-file by name and with args
func (sr *ScriptRunner) Run(name string, args []string, ctx *RunCtx) error {
	found, err := sr.TryRun(name, args, ctx)
	if !found {
		return errorx.Rawf("script file %q is not exists", name)
	}
	return err
}

// TryRun script or script-file by name and with args
func (sr *ScriptRunner) TryRun(name string, args []string, ctx *RunCtx) (found bool, err error) {
	if err := sr.InitLoad(); err != nil {
		return false, err
	}

	found = true
	ctx = EnsureCtx(ctx).WithName(name)

	// try check is script and run it.
	si, err := sr.ScriptItem(name)
	if err != nil {
		return found, err
	}
	if si != nil {
		return found, sr.doExecScript(si, args, ctx)
	}

	// try check and run script file.
	si, err = sr.ScriptFileItem(name)
	if err != nil {
		return found, err
	}

	if si != nil {
		return found, sr.doExecScriptFile(si, args, ctx)
	}
	return false, nil
}

// RunDefinedScript by input name and with arguments
func (sr *ScriptRunner) RunDefinedScript(name string, args []string, ctx *RunCtx) error {
	if err := sr.InitLoad(); err != nil {
		return err
	}

	si, err := sr.ScriptItem(name)
	if err != nil {
		return err
	}

	if si != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return sr.doExecScript(si, args, ctx)
	}
	return errorx.Rawf("script %q is not exists", name)
}

// RunScriptFile by input name and with arguments
func (sr *ScriptRunner) RunScriptFile(name string, args []string, ctx *RunCtx) error {
	if err := sr.InitLoad(); err != nil {
		return err
	}

	si, err := sr.ScriptFileItem(name)
	if err != nil {
		return err
	}

	if si != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return sr.doExecScriptFile(si, args, ctx)
	}
	return errorx.Rawf("script file %q is not exists", name)
}

func (sr *ScriptRunner) doExecScriptFile(si *ScriptItem, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si)
	}

	// run script file
	return cmdr.NewCmd(si.Bin, si.File).
		WorkDirOnNE(si.Workdir).
		WithDryRun(ctx.DryRun).
		AppendEnv(si.Env).
		AddArgs(inArgs).
		PrintCmdline().
		FlushRun()
}

func (sr *ScriptRunner) doExecScript(si *ScriptItem, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si)
	}

	ln := len(si.Cmds)
	if ln == 0 {
		return errorx.Rawf("empty cmd config for script %q", ctx.Name)
	}

	if len(inArgs) < len(si.Args) {
		return errorx.Rawf("missing required args for run script %q", ctx.Name)
	}

	// only one
	if ln == 1 {
		line := sr.handleCmdline(si.Cmds[0], inArgs, si)
		shell := strutil.OrElse(si.Type, ctx.Type)

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		return cmd.WorkDirOnNE(si.Workdir).
			WithDryRun(ctx.DryRun).
			AppendEnv(si.Env).
			PrintCmdline().
			FlushRun()
	}

	// multi command
	cr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
		rr.Workdir = si.Workdir
		rr.EnvMap = si.Env
		rr.OutToStd = true
		rr.DryRun = ctx.DryRun
	})

	shell := strutil.OrElse(si.Type, ctx.Type)
	for _, line := range si.Cmds {
		line = sr.handleCmdline(line, inArgs, si)

		if shell != "" {
			cr.CmdWithArgs(shell, "-c", line)
		} else {
			cr.AddCmdline(line)
		}
	}
	return cr.Run()
}

// process vars and env
func (sr *ScriptRunner) handleCmdline(line string, args []string, si *ScriptItem) string {
	argStr := strings.Join(args, " ")
	vars := map[string]string{
		"$@": argStr,                // 是一个字符串参数数组
		"$*": strutil.Quote(argStr), // 把所有参数合并成一个字符串
		// context info
		"$workdir": si.Workdir,
	}

	for i, val := range args {
		key := "$" + mathutil.String(i+1)
		vars[key] = val
	}

	line = strutil.Replaces(line, vars)

	// eg: $SHELL
	if sr.ParseEnv && strutil.ContainsByte(line, '$') {
		envs := sysutil.EnvironWith(si.Env)
		return textutil.RenderSMap(line, envs, "$,")
	}

	return line
}

// DefinedScript info get
func (sr *ScriptRunner) DefinedScript(name string) (any, bool) {
	info, ok := sr.Scripts[name]
	return info, ok
}

// ScriptItem get script info as ScriptItem
func (sr *ScriptRunner) ScriptItem(name string) (*ScriptItem, error) {
	info, ok := sr.Scripts[name]
	if !ok {
		// return nil, errorx.Rawf("script %q is not exists", name)
		return nil, nil
	}

	return sr.newDefinedScriptItem(name, info)
}

func (sr *ScriptRunner) newDefinedScriptItem(name string, info any) (*ScriptItem, error) {
	si := &ScriptItem{Name: name}

	switch typVal := info.(type) {
	case string: // on command
		si.Cmds = []string{typVal}
	case []any: // as commands
		si.Cmds = arrutil.SliceToStrings(typVal)
	case map[string]any: // as structured
		data := maputil.Data(typVal)
		si.Type = data.Str("type")
		si.Desc = data.Str("desc")
		si.Workdir = data.Str("workdir")

		err := si.loadArgsDefine(data.Get("args"))
		if err != nil {
			return nil, err
		}

		si.Cmds = data.Strings("cmds")
		si.Env = data.StringMap("env")
	default:
		return nil, errorx.Rawf("invalid config of the script %q", name)
	}

	if sr.WrapShell != "" && si.Type == "" {
		si.Type = sr.WrapShell
	}
	return si, nil
}

// ScriptFileItem info get
func (sr *ScriptRunner) ScriptFileItem(name string) (*ScriptItem, error) {
	// with ext
	if inExt := fsutil.FileExt(name); len(inExt) > 0 {
		fpath, ok := sr.scriptFiles[name]
		if !ok {
			return nil, nil
		}

		return sr.newFileScriptItem(name, fpath, inExt)
	}

	// auto check ext
	for _, ext := range sr.AllowedExt {
		fpath, ok := sr.scriptFiles[name+ext]
		if !ok {
			continue
		}

		return sr.newFileScriptItem(name, fpath, ext)
	}

	// not found
	return nil, nil
}

func (sr *ScriptRunner) newFileScriptItem(name, fpath, ext string) (*ScriptItem, error) {
	si := &ScriptItem{
		Name: name,
		File: fpath,
		Ext:  ext,
		Bin:  ext[1:],
	}

	if bin, ok := sr.ExtToBinMap[ext]; ok {
		si.Bin = bin
	}

	return si, nil
}

// IsDefinedScript name
func (sr *ScriptRunner) IsDefinedScript(name string) bool {
	_, ok := sr.Scripts[name]
	return ok
}

// DefinedScripts map
func (sr *ScriptRunner) DefinedScripts() map[string]any {
	return sr.Scripts
}

// ScriptFiles file map
func (sr *ScriptRunner) ScriptFiles() map[string]string {
	return sr.scriptFiles
}

func (sr *ScriptRunner) WithConfigFn(fn func(sr *ScriptRunner)) *ScriptRunner {
	fn(sr)
	return sr
}
