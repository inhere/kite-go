package kscript

import (
	"os"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// Runner struct
type Runner struct {
	// Scripts config and loaded from DefineFiles.
	//
	// format: {name: info, name2: info2, ...}
	Scripts map[string]any `json:"scripts"`

	// DefineDir scripts define dir, will read and add to Scripts
	DefineDir string `json:"define_dir"`

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
func (r *Runner) Search(name string, args []string, limit int) map[string]string {
	result := make(map[string]string)
	limit = mathutil.Min(limit, 3)
	goutil.MustOK(r.InitLoad())

	// TODO use args for limit search

	for sName, sInfo := range r.Scripts {
		if strutil.IContains(sName, name) {
			result[sName] = strutil.Truncate(goutil.String(sInfo), 48, "...")
			if limit > 0 && len(result) >= limit {
				return result
			}
		}
	}

	// search script files
	for fName, fPath := range r.scriptFiles {
		if strutil.IContains(fName, name) {
			result[fName] = fPath
			if limit > 0 && len(result) >= limit {
				return result
			}
		}
	}

	return result
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
		return found, r.runDefineScript(si, args, ctx)
	}

	// try check and run script file.
	si, err = r.ScriptFileInfo(name)
	if err != nil {
		return found, err
	}

	if si != nil {
		return found, r.runScriptFile(si, args, ctx)
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
		return r.runDefineScript(si, args, ctx)
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
		return r.runScriptFile(si, args, ctx)
	}
	return errorx.Rawf("script file %q is not exists", name)
}

// RunScriptInfo by args and context
func (r *Runner) RunScriptInfo(si *ScriptInfo, inArgs []string, ctx *RunCtx) error {
	if si.IsFile() {
		return r.runScriptFile(si, inArgs, ctx)
	}
	return r.runDefineScript(si, inArgs, ctx)
}

func (r *Runner) runScriptFile(si *ScriptInfo, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si, ctx)
	}

	// run script file
	return cmdr.NewCmd(si.BinName, si.File).
		WorkDirOnNE(si.Workdir).
		WithDryRun(ctx.DryRun).
		AppendEnv(si.Env).
		AddArgs(inArgs).
		PrintCmdline().
		FlushRun()
}

func (r *Runner) runDefineScript(si *ScriptInfo, inArgs []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(si, ctx)
	}

	ln := len(si.Cmds)
	if ln == 0 {
		return errorx.Rawf("empty cmd config for script %q", ctx.Name)
	}

	needArgs := si.ParseArgs()
	if nln := len(needArgs); len(inArgs) < nln {
		return errorx.Rawf("missing required args for run script %q(need %d)", ctx.Name, nln)
	}

	envMap := ctx.MergeEnv(si.Env)
	shell := strutil.OrElse(ctx.Type, si.Type)
	workdir := strutil.OrElse(ctx.Workdir, si.Workdir)

	// build context vars
	argStr := strings.Join(inArgs, " ")
	vars := map[string]string{
		// 是一个字符串参数数组
		"$@": argStr,
		"$*": strutil.Quote(argStr), // 把所有参数合并成一个字符串
		// context info
		"$workdir": workdir,
		"$dirname": fsutil.Name(workdir),
	}

	for i, val := range inArgs {
		key := "$" + mathutil.String(i+1)
		vars[key] = val
	}

	// exec each command
	for _, line := range si.Cmds {
		if len(line) == 0 {
			continue
		}

		// redirect run other script
		if line[0] == '@' {
			name := line[1:]
			osi, err := r.ScriptDefineInfo(name)
			if err != nil {
				return err
			}
			if osi == nil {
				return errorx.Rawf("run %q: reference script %q not found", si.Name, name)
			}

			err = r.runDefineScript(osi, inArgs, ctx)
			if err != nil {
				return err
			}
			continue
		}

		line = r.handleCmdline(line, vars, si)

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		err := cmd.WorkDirOnNE(workdir).
			WithDryRun(ctx.DryRun).
			AppendEnv(envMap).
			PrintCmdline().
			FlushRun()

		if err != nil {
			return err
		}
	}
	return nil
}

// process vars and env
func (r *Runner) handleCmdline(line string, vars map[string]string, si *ScriptInfo) string {
	line = strutil.Replaces(line, vars)

	// eg: $SHELL
	if r.ParseEnv && strutil.ContainsByte(line, '$') {
		envs := sysutil.EnvMapWith(si.Env)
		return textutil.RenderSMap(line, envs, "$,")
	}

	return line
}

// RawDefinedScript raw info get
func (r *Runner) RawDefinedScript(name string) (any, bool) {
	info, ok := r.Scripts[name]
	return info, ok
}

// ScriptDefineInfo get script info as ScriptInfo
func (r *Runner) ScriptDefineInfo(name string) (*ScriptInfo, error) {
	info, ok := r.Scripts[name]
	if !ok {
		return nil, nil // not found TODO ErrNotFund
	}
	return newDefinedScriptInfo(name, info, r.TypeShell)
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
		Name:    name,
		File:    fpath,
		FileExt: ext,
		BinName: ext[1:],
	}

	if bin, ok := r.ExtToBinMap[ext]; ok {
		si.BinName = bin
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
