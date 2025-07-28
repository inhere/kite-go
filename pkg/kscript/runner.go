package kscript

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/gookit/slog"
)

var settingsKey = "__settings"

type RunnerMeta struct {
	// apps  []*ScriptApp
	// tasks []*ScriptTask
	// files []*ScriptFile
}

// Runner struct. TODO KRunner, ScriptRunner or ScriptManager
//
// 实现扩展的kite run命令，可以执行任何的 script-file, script-task, script-app 等等
type Runner struct {
	RunnerMeta

	// PathResolver handler. 用于查找脚本文件
	PathResolver func(path string) string

	// ------------------------ config for script app --------------------

	// ScriptAppDirs 独立的 script app 定义文件目录. 每个定义文件是一个独立的cli app
	//  - default will load from `$base/script-app`
	ScriptAppDirs []string `json:"script_app_dirs"`
	// ScriptAppExts script app file define extensions. eg: .yml, .yaml
	ScriptAppExts []string `json:"script_app_exts"`

	// 加载并解析后的 app 定义
	apps map[string]*ScriptApp
	// files loaded from ScriptAppDirs. format: {filename: filepath, ...}
	appFiles map[string]string
	// mark script app loaded
	appLoaded bool

	// ------------------------ config for script task --------------------

	// DefineFiles script tasks define files, will read and add to Scripts
	//
	// Allow vars: $user, $os
	//
	// eg:
	//	- config/module/scripts.yml
	//	- ?config/module/scripts.$os.yml  // start withs '?' - an optional file, load on exists.
	DefineFiles []string `json:"define_files"`
	// 自动加载的task文件名称列表，无需设置扩展
	//  - 将自动从当前目录或父级目录中寻找 script task 定义文件
	//  - 找到第一个匹配的就停止
	AutoTaskFiles []string `json:"auto_task_files"`
	// 自动加载的task文件扩展名
	AutoTaskExts []string `json:"auto_task_exts"`
	// auto 向上搜索目录最大深度，默认为 6. 找到第一个匹配的就停止
	AutoMaxDepth int `json:"auto_max_depth"`

	// Scripts 通过配置定义的各种简单的任务命令。tasks config and loaded from DefineFiles.
	//
	// Format: {name: info, name2: info2, ...}
	//
	//  - special settings key: __settings, will read and merge to Settings
	Scripts map[string]any `json:"scripts"`

	// ParseEnv var on script command expr. eg: $SHELL
	ParseEnv bool `json:"parse_env"`
	// TypeShell wrapper for run each script.
	//
	// value like: bash, sh, zsh, cmd, pwsh or empty for direct run command
	TypeShell string `json:"type_shell"`

	// 加载并解析后的 tasks 定义
	// tasks *ScriptTasks TODO
	tasks map[string]*ScriptTask
	// mark script task loaded
	taskLoaded bool
	// settings for all script tasks.
	//
	// eg:
	//  - vars: map[string]string built in vars map. group name: vars
	//  - group: map[string]map[string]string grouped var map.
	taskSettings TaskSettings

	// ------------------------ config for script file --------------------

	// ScriptDirs 独立的 script file 文件查找目录。例如 bash, python, php 等脚本文件
	ScriptDirs []string `json:"script_dirs"`

	// AllowedExt allowed script file extensions. eg: .go, .sh
	AllowedExt []string `json:"allowed_ext"`
	// FindBinByExt on run a script file
	FindBinByExt bool `json:"find_bin_by_ext"`
	// ExtToBinMap settings. key: ext, value: bin name or path
	ExtToBinMap map[string]string `json:"ext_to_bin_map"`
	// BinPathMap settings. key: bin name, value: bin path
	BinPathMap map[string]string `json:"bin_path_map"`

	// loaded from ScriptDirs. format: {filename: filepath, ...}
	scriptFiles map[string]string
	fileMetas map[string]*ScriptFile
	// mark script loaded
	fileLoaded bool
}

/*
-----------
--------------------------------- Init load ---------------------------------
----------- region T: Init load
*/

// InitLoad define scripts and script files.
func (r *Runner) InitLoad() error {
	if err := r.LoadScriptTaskInfos(); err != nil {
		return err
	}

	r.LoadScriptApps()

	return r.LoadScriptFiles()
}

/* endregion
--------------------------------- Load task files ---------------------------------
----------- region T: Load task files
*/

// LoadScriptTaskInfos from Runner.DefineFiles
func (r *Runner) LoadScriptTaskInfos() (err error) {
	if r.taskLoaded {
		return nil
	}

	r.taskLoaded = true
	loader := config.New("loader")
	loader.AddDriver(ini.Driver)
	loader.AddDriver(yaml.Driver)
	loader.AddDriver(toml.Driver)

	// 从配置的定义文件中加载
	for _, fPath := range r.DefineFiles {
		// optional file
		var optional bool
		if fPath[0] == '?' {
			optional = true
			fPath = fPath[1:]
		}

		fPath = r.PathResolver(fPath)
		if optional && !fsutil.IsFile(fPath) {
			continue
		}

		slog.Debugf("load script task file %q", fPath)
		err = loader.LoadFiles(fPath)
		if err != nil {
			return errorx.Errorf("load task file %q error: %s", fPath, err)
		}

		r.Scripts = maputil.SimpleMerge(loader.Data(), r.Scripts)
		loader.ClearData()
	}

	// 从工作目录/父级目录自动加载
	if fPaths := r.findAutoTaskFiles(); len(fPaths) > 0 {
		for _, fPath := range fPaths {
			err = loader.LoadFiles(fPath)
			if err != nil {
				return errorx.Wrapf(err, "load auto task file %q error: %s", fPath, err)
			}

			r.Scripts = maputil.SimpleMerge(loader.Data(), r.Scripts)
			loader.ClearData()
		}
	}

	// load custom settings
	if setData, ok := r.Scripts[settingsKey]; ok {
		delete(r.Scripts, settingsKey)
		if setMap, ok1 := setData.(map[string]any); ok1 {
			r.taskSettings.loadData(setMap)
		}
	}

	return nil
}

// 从工作目录/父级目录自动查找 task 定义文件,向上层级越高的文件在前面(先加载)
func (r *Runner) findAutoTaskFiles() (ss []string) {
	findDir := sysutil.Workdir()
	findLevel := 1

	// 从当前目录或父级目录中寻找 script task 配置文件
	for {
		// 一个目录下只匹配一个文件，找到一个就停止。
		var founded bool

		for _, fName := range r.AutoTaskFiles {
			for _, ext := range r.AutoTaskExts {
				fPath := findDir + "/" + fName + ext
				if fsutil.IsFile(fPath) {
					slog.Debugf("found task file %q", fPath)
					ss = append(ss, fPath)
					founded = true
					break
				}
			}
			if founded {
				break
			}
		}

		if findLevel >= r.AutoMaxDepth {
			break
		}

		findLevel++
		findDir = filepath.Dir(findDir)
		if len(findDir) < 3 {
			break
		}
	}

	// 倒序, 从最顶层开始
	if len(ss) > 0 {
		arrutil.Reverse(ss)
	}
	return
}

/* endregion
--------------------------------- Load script apps ---------------------------------
----------- region T: Load script apps
*/

// LoadScriptApps from Runner.ScriptApps
func (r *Runner) LoadScriptApps() {
	if r.appLoaded {
		return
	}
	r.appLoaded = true

	for _, dirPath := range r.ScriptAppDirs {
		dirPath = r.PathResolver(dirPath)
		des, err := os.ReadDir(dirPath)
		if err != nil {
			slog.Warnf("kscript: read dir %q error: %s", dirPath, err)
			continue
		}

		for _, ent := range des {
			fName := ent.Name()
			if !ent.IsDir() {
				nameNoExt := fsutil.NameNoExt(fName)
				fullPath := dirPath + "/" + fName
				r.appFiles[nameNoExt] = fullPath
				slog.Debugf("kscript: load script app %q(path: %s)", nameNoExt, fullPath)
			}
		}
	}
}

/* endregion
--------------------------------- Load script files ---------------------------------
----------- region T: Load script files
*/

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
			slog.Warnf("kscript: read dir %q error: %s", dirPath, err)
			continue
		}

		for _, ent := range des {
			fName := ent.Name()
			if !ent.IsDir() {
				fullPath := dirPath + "/" + fName
				r.scriptFiles[fName] = fullPath
				slog.Debugf("kscript: load script file %q(path: %s)", fName, fullPath)
			}
		}
	}

	return nil
}

/* endregion
------------------------------------------------------------------
----------- region T: Search script
*/

// Search1ByName search one script task/file by name
func (r *Runner) Search1ByName(name string, limit int) string {
	ret := r.SearchByName(name, 1)
	if len(ret) > 0 {
		for na := range ret {
			return na
		}
	}
	return ""
}

// SearchByName search script task/file by name
func (r *Runner) SearchByName(name string, limit int) map[string]string {
	parts := strutil.SplitTrimmed(name, " ")
	return r.SearchByKeywords(parts, limit)
}

// SearchByKeywords search script task/file by keywords
func (r *Runner) SearchByKeywords(parts []string, limit int) map[string]string {
	ret := map[string]string{}

	// find in script tasks
	for sName, sInfo := range r.Scripts {
		if strutil.IContainsAll(sName, parts) {
			ret[sName] = strutil.Truncate(goutil.String(sInfo), 68, "...")
			if limit > 0 && len(ret) >= limit {
				return ret
			}
		}
	}

	// search script files
	for fName, fPath := range r.scriptFiles {
		if strutil.IContainsAll(fName, parts) {
			ret[fName] = fPath
			if limit > 0 && len(ret) >= limit {
				return ret
			}
		}
	}

	return ret
}

// Search by name or description
func (r *Runner) Search(name string, args []string, limit int) map[string]string {
	result := make(map[string]string)
	limit = mathutil.Min(limit, 3)
	goutil.MustOK(r.InitLoad())

	parts := []string{name}
	if strutil.ContainsByte(name, ' ') {
		parts = strutil.SplitTrimmed(name, " ")
	}
	// append args to parts
	// TODO use args for limit search
	parts = append(parts, args...)

	for sName, sInfo := range r.Scripts {
		if strutil.IContainsAll(sName, parts) {
			result[sName] = strutil.Truncate(goutil.String(sInfo), 68, "...")
			if limit > 0 && len(result) >= limit {
				return result
			}
		}
	}

	// search script files
	for fName, fPath := range r.scriptFiles {
		if strutil.IContainsAll(fName, parts) {
			result[fName] = fPath
			if limit > 0 && len(result) >= limit {
				return result
			}
		}
	}

	return result
}

/*
----------- endregion
--------------------------------- Run script ---------------------------------
----------- region T: Run script
*/

// Run script or script-file by name and with args
func (r *Runner) Run(name string, args []string, ctx *RunCtx) error {
	found, err := r.TryRun(name, args, ctx)
	if !found {
		return errorx.Rawf("script file %q is not exists", name)
	}
	return err
}

// TryRun script task or script-file by name and with args
func (r *Runner) TryRun(name string, args []string, ctx *RunCtx) (found bool, err error) {
	if err := r.InitLoad(); err != nil {
		return false, err
	}

	found = true
	ctx = EnsureCtx(ctx).WithName(name)

	// ------ try check is task and run it ------
	si, err := r.LoadScriptTaskInfo(name)
	if err != nil {
		return found, err
	}
	if si != nil {
		ccolor.Magentaln("Run script task:", name)
		return found, r.runScriptTask(si, args, ctx)
	}

	// ------ try check is file and run it ------
	sf, err := r.LoadScriptFileInfo(name)
	if err != nil {
		return found, err
	}

	if sf != nil {
		ccolor.Magentaln("Run script file: %s", name)
		return found, r.runScriptFile(sf, args, ctx)
	}
	return false, nil
}

/*
----------- endregion
--------------------------------- Run script task ---------------------------------
----------- region T: Run script task
*/

// RunScriptTask by input name and with arguments
func (r *Runner) RunScriptTask(name string, args []string, ctx *RunCtx) error {
	if err := r.InitLoad(); err != nil {
		return err
	}

	si, err := r.LoadScriptTaskInfo(name)
	if err != nil {
		return err
	}

	if si != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return r.runScriptTask(si, args, ctx)
	}
	return errorx.Rawf("script task %q is not exists", name)
}

func (r *Runner) runScriptTask(st *ScriptTask, inArgs []string, ctx *RunCtx) error {
	ctx.ScriptType = TypeTask
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(st, ctx)
	}

	ln := len(st.Cmds)
	if ln == 0 {
		return errorx.Rawf("empty cmd config for script %q", ctx.Name)
	}

	needArgs := st.ParseArgs()
	if nln := len(needArgs); len(inArgs) < nln {
		ccolor.Println("<mga>Script task contents:</>\n ", strings.Join(st.Cmds, "\n  "))
		return errorx.Rawf("missing required args for run task %q(need %d)", ctx.Name, nln)
	}

	envMap := ctx.MergeEnv(st.Env)
	shell := strutil.OrElse(ctx.Type, st.Type)
	workdir := strutil.OrElse(ctx.Workdir, st.Workdir)

	// build context vars
	argStr := strings.Join(inArgs, " ")
	vars := map[string]any{
		// $@ 是一个字符串参数数组
		"@": argStr,
		// @* 把所有参数合并成一个字符串
		"*": strutil.Quote(argStr),
		// context info
		"workdir": workdir,
		"dirname": fsutil.Name(workdir),
	}

	// 输入参数处理 $1 ... $N
	for i, val := range inArgs {
		// key := "$" + mathutil.String(i+1)
		key := mathutil.String(i + 1)
		vars[key] = val
	}

	// 追加 全局变量
	vars = r.buildTaskTplVars(vars)
	vars["ctx"] = ctx.Vars
	if ctx.AppendVarsFn != nil {
		vars = ctx.AppendVarsFn(vars)
	}
	if ctx.Verbose {
		show.AList("Task Vars", vars)
	}
	if workdir != "" {
		ccolor.Magentaln("Workdir:", workdir)
	}

	// exec each command
	for _, line := range st.Cmds {
		if len(line) == 0 {
			continue
		}

		// redirect run another script
		if line[0] == '@' {
			name := line[1:]
			osi, err := r.LoadScriptTaskInfo(name)
			if err != nil {
				return err
			}
			if osi == nil {
				return errorx.Rawf("run %q: reference script %q not found", st.Name, name)
			}

			err = r.runScriptTask(osi, inArgs, ctx)
			if err != nil {
				return err
			}
			continue
		}

		line = r.handleCmdline(line, vars, st)

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		err := cmd.WorkDirOnNE(workdir).WithDryRun(ctx.DryRun).AppendEnv(envMap).PrintCmdline2().FlushRun()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) buildTaskTplVars(data map[string]any) map[string]any {
	tn := time.Now()
	data["time"] = map[string]any{
		"unix_sec":    tn.Unix(),
		"datetime":    tn.Format("2006-01-02 15:04:05"),
		"date_Ymd_hm": tn.Format("2006-01-02_15:04"),
		"date_ymd_hm": tn.Format("06-01-02_15:04"),
		"date_ymd":    tn.Format("2006-01-02"),
		"date_hms":    tn.Format("15:04:05"),
	}

	// vars in runner.taskSettings
	data["vars"] = r.taskSettings.Vars
	data["groups"] = r.taskSettings.Groups

	return data
}

// 使用简单的模板渲染，支持链式语法变量替换，环境变量，默认值等 - 无法同时支持 $var ${var_name}
// var rpl = textutil.NewVarReplacer("$").WithParseEnv().WithParseDefault()
// 专门实现的类似 php, shell 的字符串表达式处理
var rpl = textutil.NewStrVarRenderer()

// process vars and env
func (r *Runner) handleCmdline(line string, vars map[string]any, st *ScriptTask) string {
	envs := sysutil.EnvMapWith(st.Env)

	rpl.SetGetter(func(name string) (val string, ok bool) {
		// eg: $SHELL
		if r.ParseEnv {
			if val, ok = envs[name]; ok {
				return val, true
			}
		}
		return name, false
	})

	// line = strutil.Replaces(line, vars)
	// eg: $SHELL
	// if r.ParseEnv && strutil.ContainsByte(line, '$') {
	// 	envs := sysutil.EnvMapWith(st.Env)
	// 	return textutil.RenderSMap(line, envs, "$,")
	// }

	return rpl.Render(line, vars)
}

// RawScriptTask raw info get
func (r *Runner) RawScriptTask(name string) (any, bool) {
	info, ok := r.Scripts[name]
	return info, ok
}

// LoadScriptTaskInfo get script info as ScriptTask
func (r *Runner) LoadScriptTaskInfo(name string) (*ScriptTask, error) {
	// TODO 先读取 Runner.tasks 缓存，如果找不到再从 Scripts 中解析读取

	info, ok := r.Scripts[name]
	if !ok {
		return nil, nil // not found TODO ErrNotFound
	}
	return parseScriptTask(name, info, r.TypeShell)
}

/*
----------- endregion
-----------------------------------------------------------------------------
----------- region T: Run script file
*/

// RunScriptFile by input name and with arguments
func (r *Runner) RunScriptFile(name string, args []string, ctx *RunCtx) error {
	if err := r.InitLoad(); err != nil {
		return err
	}

	sf, err := r.LoadScriptFileInfo(name)
	if err != nil {
		return err
	}

	if sf != nil {
		ctx = EnsureCtx(ctx).WithName(name)
		return r.runScriptFile(sf, args, ctx)
	}
	return errorx.Rawf("script file %q is not exists", name)
}

func (r *Runner) runScriptFile(sf *ScriptFile, inArgs []string, ctx *RunCtx) error {
	ctx.ScriptType = TypeFile
	if ctx.BeforeFn != nil {
		ctx.BeforeFn(sf, ctx)
	}

	// run script file
	return cmdr.NewCmd(sf.BinName, sf.File).
		WorkDirOnNE(sf.Workdir).
		WithDryRun(ctx.DryRun).
		AppendEnv(sf.Env).
		AddArgs(inArgs).
		PrintCmdline2().
		FlushRun()
}

// LoadScriptFileInfo info get
func (r *Runner) LoadScriptFileInfo(name string) (*ScriptFile, error) {
	// with ext
	if inExt := fsutil.FileExt(name); len(inExt) > 0 {
		fPath, ok := r.scriptFiles[name]
		if !ok {
			return nil, nil
		}

		return r.newScriptFileInfo(name, fPath, inExt)
	}

	// auto check ext
	for _, ext := range r.AllowedExt {
		fPath, ok := r.scriptFiles[name+ext]
		if !ok {
			continue
		}

		return r.newScriptFileInfo(name, fPath, ext)
	}

	// not found
	return nil, nil
}

func (r *Runner) newScriptFileInfo(name, fPath, ext string) (*ScriptFile, error) {
	si := &ScriptFile{
		ScriptMeta: ScriptMeta{
			ScriptType: TypeFile,
		},
		Name:    name,
		File: fPath,
		FileExt: ext,
		BinName: ext[1:],
	}

	if bin, ok := r.ExtToBinMap[ext]; ok {
		si.BinName = bin
	}
	return si, nil
}

// IsScriptTask name
func (r *Runner) IsScriptTask(name string) bool {
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
