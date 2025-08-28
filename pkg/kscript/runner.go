package kscript

import (
	"os"
	"path/filepath"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
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
	if err := r.LoadScriptTasks(); err != nil {
		return err
	}

	// load script apps
	r.LoadScriptApps()

	// load script files
	return r.LoadScriptFiles()
}

/* endregion
--------------------------------- Load task files ---------------------------------
----------- region T: Load task files
*/

// LoadScriptTasks from Runner.DefineFiles
func (r *Runner) LoadScriptTasks() (err error) {
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

// LoadScriptTaskInfo get script info as ScriptTask
func (r *Runner) LoadScriptTaskInfo(name string) (*ScriptTask, error) {
	// TODO 先读取 Runner.tasks 缓存，如果找不到再从 Scripts 中解析读取

	info, ok := r.Scripts[name]
	if !ok {
		return nil, nil // not found TODO ErrNotFound
	}

	// TODO 支持别名查找

	return parseScriptTask(name, info, r.TypeShell)
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

/* endregion
------------------------------------------------------------------
----------- region T: helper methods
*/

// RawScriptTask raw info get
func (r *Runner) RawScriptTask(name string) (any, bool) {
	info, ok := r.Scripts[name]
	return info, ok
}

// IsScriptTask name
func (r *Runner) IsScriptTask(name string) bool {
	_, ok := r.Scripts[name]
	return ok
}

// RawScriptTasks map
func (r *Runner) RawScriptTasks() map[string]any {
	return r.Scripts
}

// ScriptFiles file map
func (r *Runner) ScriptFiles() map[string]string {
	return r.scriptFiles
}
