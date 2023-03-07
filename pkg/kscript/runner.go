package kscript

import (
	"os"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
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

// Runner struct
type Runner struct {
	// Scripts config and loaded from DefineFiles.
	//
	// format: {name: info, name2: info2, ...}
	Scripts map[string]any `json:"scripts"`
	// DefineFiles scripts define files, will read and add to Scripts
	//
	// Allow vars: $user, $os
	//
	// eg:
	//	- config/module/scripts.yml
	//	- ?config/module/scripts.$os.yml  // start withs '?' - an optional file, load on exists.
	DefineFiles []string `json:"define_files"`
	// TypeShell wrapper for run each script.
	//
	// value like: bash, sh, zsh or empty for direct run command
	TypeShell string `json:"type_shell"`
	// ParseEnv var on script command
	ParseEnv bool `json:"parse_env"`
	// Aliases script name aliases map
	Aliases map[string]string `json:"aliases"`

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

/*
-----------
--------------------------------- Init load ---------------------------------
-----------
*/

// InitLoad define scripts and script files.
func (r *Runner) InitLoad() error {
	if len(r.Aliases) > 0 {
		// TODO format
	}

	if err := r.LoadDefineScripts(); err != nil {
		return err
	}
	return r.LoadScriptFiles()
}

// LoadDefineScripts from DefineFiles
func (r *Runner) LoadDefineScripts() (err error) {
	if r.defineLoaded {
		return nil
	}

	r.defineLoaded = true
	loader := config.New("loader")
	loader.AddDriver(ini.Driver)
	loader.AddDriver(yaml.Driver)
	loader.AddDriver(toml.Driver)

	for _, fpath := range r.DefineFiles {
		// optional file
		var optional bool
		if fpath[0] == '?' {
			optional = true
			fpath = fpath[1:]
		}

		fpath = r.PathResolver(fpath)
		if optional {
			err = loader.LoadExists(fpath)
		} else {
			err = loader.LoadFiles(fpath)
		}

		if err != nil {
			return err
		}

		r.Scripts = maputil.SimpleMerge(loader.Data(), r.Scripts)
		loader.ClearData()
	}

	return nil
}

// LoadScriptFiles from the ScriptDirs
func (r *Runner) LoadScriptFiles() error {
	if r.fileLoaded {
		return nil
	}
	r.fileLoaded = true

	for _, dirPath := range r.ScriptDirs {
		dirPath = r.PathResolver(dirPath)
		des, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, ent := range des {
			fName := ent.Name()
			if !ent.IsDir() {
				// add
				r.scriptFiles[fName] = dirPath + "/" + fName
			}
		}
	}

	return nil
}

// Search by name
func (r *Runner) Search(name string, ctx *RunCtx) error {
	return nil
}

/*
-----------
--------------------------------- Run script ---------------------------------
-----------
*/

// Run script or script-file by name and with args
func (r *Runner) Run(name string, args []string, ctx *RunCtx) error {
	found, err := r.TryRun(name, args, ctx)
	if !found {
		return errorx.Rawf("script file %q is not exists", name)
	}
	return err
}

// TryRun script or script-file by name and with args
func (r *Runner) TryRun(name string, args []string, ctx *RunCtx) (found bool, err error) {
	if err := r.InitLoad(); err != nil {
		return false, err
	}

	found = true
	ctx = EnsureCtx(ctx).WithName(name)

	// try check is script and run it.
	si, err := r.ScriptDefineInfo(name)
	if err != nil {
		return found, err
	}
	if si != nil {
		return found, r.doExecScript(si, args, ctx)
	}

	// try check and run script file.
	si, err = r.ScriptFileInfo(name)
	if err != nil {
		return found, err
	}

	if si != nil {
		return found, r.doExecScriptFile(si, args, ctx)
	}
	return false, nil
}

// RunDefinedScript by input name and with arguments
func (r *Runner) RunDefinedScript(name string, args []string, ctx *RunCtx) error {
	if err := r.InitLoad(); err != nil {
		return err
	}

	si, err := r.ScriptDefineInfo(name)
	if err != nil {
		return err
	}

	if si != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return r.doExecScript(si, args, ctx)
	}
	return errorx.Rawf("script %q is not exists", name)
}

// RunScriptFile by input name and with arguments
func (r *Runner) RunScriptFile(name string, args []string, ctx *RunCtx) error {
	if err := r.InitLoad(); err != nil {
		return err
	}

	si, err := r.ScriptFileInfo(name)
	if err != nil {
		return err
	}

	if si != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return r.doExecScriptFile(si, args, ctx)
	}
	return errorx.Rawf("script file %q is not exists", name)
}

func (r *Runner) doExecScriptFile(si *ScriptInfo, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si, ctx)
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

func (r *Runner) doExecScript(si *ScriptInfo, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si, ctx)
	}

	ln := len(si.Cmds)
	if ln == 0 {
		return errorx.Rawf("empty cmd config for script %q", ctx.Name)
	}

	if len(inArgs) < len(si.Args) {
		return errorx.Rawf("missing required args for run script %q", ctx.Name)
	}

	shell := strutil.OrElse(ctx.Type, si.Type)
	workdir := strutil.OrElse(ctx.Workdir, si.Workdir)

	// only one
	if ln == 1 {
		line := r.handleCmdline(si.Cmds[0], inArgs, si)

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		return cmd.WorkDirOnNE(workdir).
			WithDryRun(ctx.DryRun).
			AppendEnv(ctx.MergeEnv(si.Env)).
			PrintCmdline().
			FlushRun()
	}

	// multi command
	cr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
		rr.OutToStd = true
		rr.Workdir = workdir
		rr.EnvMap = ctx.MergeEnv(si.Env)
		rr.DryRun = ctx.DryRun
	})

	for _, line := range si.Cmds {
		line = r.handleCmdline(line, inArgs, si)

		if shell != "" {
			cr.CmdWithArgs(shell, "-c", line)
		} else {
			cr.AddCmdline(line)
		}
	}
	return cr.Run()
}

// process vars and env
func (r *Runner) handleCmdline(line string, args []string, si *ScriptInfo) string {
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
	if r.ParseEnv && strutil.ContainsByte(line, '$') {
		envs := sysutil.EnvironWith(si.Env)
		return textutil.RenderSMap(line, envs, "$,")
	}

	return line
}

// DefinedScript info get
func (r *Runner) DefinedScript(name string) (any, bool) {
	info, ok := r.Scripts[name]
	return info, ok
}

// ScriptDefineInfo get script info as ScriptInfo
func (r *Runner) ScriptDefineInfo(name string) (*ScriptInfo, error) {
	info, ok := r.Scripts[name]
	if !ok {
		return nil, nil // not found
	}
	return r.newDefinedScriptInfo(name, info)
}

func (r *Runner) newDefinedScriptInfo(name string, info any) (*ScriptInfo, error) {
	si := &ScriptInfo{Name: name}

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

	si.InitType(r.TypeShell)
	return si, nil
}

// ScriptFileInfo info get
func (r *Runner) ScriptFileInfo(name string) (*ScriptInfo, error) {
	// with ext
	if inExt := fsutil.FileExt(name); len(inExt) > 0 {
		fpath, ok := r.scriptFiles[name]
		if !ok {
			return nil, nil
		}

		return r.newFileScriptItem(name, fpath, inExt)
	}

	// auto check ext
	for _, ext := range r.AllowedExt {
		fpath, ok := r.scriptFiles[name+ext]
		if !ok {
			continue
		}

		return r.newFileScriptItem(name, fpath, ext)
	}

	// not found
	return nil, nil
}

func (r *Runner) newFileScriptItem(name, fpath, ext string) (*ScriptInfo, error) {
	si := &ScriptInfo{
		Name: name,
		File: fpath,
		Ext:  ext,
		Bin:  ext[1:],
	}

	if bin, ok := r.ExtToBinMap[ext]; ok {
		si.Bin = bin
	}
	return si, nil
}

// IsDefinedScript name
func (r *Runner) IsDefinedScript(name string) bool {
	_, ok := r.Scripts[name]
	return ok
}

// DefinedScripts map
func (r *Runner) DefinedScripts() map[string]any {
	return r.Scripts
}

// ScriptFiles file map
func (r *Runner) ScriptFiles() map[string]string {
	return r.scriptFiles
}
