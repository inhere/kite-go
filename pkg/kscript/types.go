package kscript

import (
	"path/filepath"
	"time"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
)

//
// endregion
// region T: script metadata
//

// ScriptType definition
type ScriptType string

const (
	TypeFile ScriptType = "file"
	TypeTask ScriptType = "task" // task define.
	TypeApp  ScriptType = "app"  // app define. 独立的cli app配置文件
)

type ScriptMeta struct {
	// ScriptType name
	ScriptType ScriptType
	// Workdir for run script, default is current dir.
	Workdir string
	// CleanEnv for run script, default is false.
	CleanEnv bool
	// Env setting for run the script file/app/task
	//  - 支持使用 ${var_name} 引用变量
	Env map[string]string
	// EnvPaths custom prepend set ENV PATH.
	//  - 支持使用 ${var_name} 引用变量
	EnvPaths []string
	// Timeout for run a script, default is 0.
	Timeout time.Duration

	parsedEnv map[string]string
}

type ScriptItem interface {
	ScriptTask | ScriptApp | ScriptFile
}

//
// endregion
// region T: script_app
//

type ScriptApps struct {
	apps map[string]*ScriptApp
}

type ScriptApp struct {
	ScriptMeta
	// script app name, use file name
	Name string
	// File script app file path in Runner.ScriptApps
	File string
	// from __task_config: actions, toplevel
	// Actions []TaskCmd
}

//
// endregion
// region T: script_file
//

type ScriptFiles struct {
	// loaded from ScriptDirs.
	//
	// format: {filename: filepath, ...}
	files map[string]string
	metas map[string]*ScriptFile
}

type ScriptFile struct {
	ScriptMeta

	// TODO read and parse file metadata.
	parsed bool

	// script name, default uses file name. eg: demo.sh
	Name string
	// File script file path in Runner.ScriptDirs
	File string
	// BinName script file bin name. 默认从 ext 解析 e.g.: .php => php
	BinName string
	// file ext. eg: .go
	FileExt string
	// ShellBang script file shell bang line.
	// always at first line and start with: #!
	ShellBang string
}

// Exec the script file with context
func (sf *ScriptFile) Exec(args []string, ctx *RunCtx) error {
	if ctx.BeforeFn != nil {
		// ctx.BeforeFn(sf, ctx)
	}

	// run script file
	return cmdr.NewCmd(sf.Name, sf.File).
		WorkDirOnNE(sf.Workdir).
		WithDryRun(ctx.DryRun).
		AppendEnv(sf.Env).
		AddArgs(args).
		PrintCmdline().
		FlushRun()
}

//
// endregion
// region T: run context
//

// RunCtx for run a task/file/app
type RunCtx struct {
	// Name for script run
	Name string
	Type string // shell type
	// ScriptType name
	ScriptType ScriptType

	// Verbose show more info on run
	Verbose bool
	// DryRun script
	DryRun bool
	// Workdir for run a script
	Workdir string
	// Vars for run cmd. access: $ctx.var_name
	Vars map[string]string
	// Env setting for run
	//  - 支持使用 ${var_name} 引用变量
	Env map[string]string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
	// Args for run script task
	Args []string

	// BeforeFn hook. si: ScriptTask | ScriptApp | ScriptFile
	BeforeFn func(si any, ctx *RunCtx)
	// AppendVarsFn hook for run task. eg: gvs, paths, kite
	AppendVarsFn func(data map[string]any) map[string]any

	fullEnv map[string]string
}

// EnsureCtx to
func EnsureCtx(ctx *RunCtx) *RunCtx {
	if ctx == nil {
		ctx = &RunCtx{}
	}

	// ensure Env
	if ctx.Env == nil {
		ctx.Env = map[string]string{}
	}
	return ctx
}

// WithName to ctx
func (c *RunCtx) WithName(name string) *RunCtx {
	c.Name = name
	return c
}

// WithNameArgs to ctx
func (c *RunCtx) WithNameArgs(name string, args []string) *RunCtx {
	c.Name = name
	c.Args = args
	return c
}

// MergeEnv and returns
func (c *RunCtx) MergeEnv(mps ...map[string]string) {
	if len(c.Env) > 0 {
		mps = append(mps, c.Env)
	}
	c.Env = maputil.MergeMultiSMap(mps...)
}

// FullEnv for run script
func (c *RunCtx) FullEnv() map[string]string {
	if len(c.fullEnv) == 0 {
		c.fullEnv = sysutil.EnvMapWith(c.Env)
	}
	return c.fullEnv
}

// 专门实现的类似 php, shell 的字符串表达式处理
var svRender = textutil.NewStrVarRenderer()

func (c *RunCtx) ParseVarInEnv(envPaths []string, vars map[string]any) map[string]string {
	// merge env
	envMap := c.Env

	// parse env expression value
	if len(envMap) > 0 && len(vars) > 0 {
		for k, v := range envMap {
			if strutil.ContainsByte(v, '$') {
				envMap[k] = svRender.Render(v, vars)
			}
		}
	}

	// merge env PATH
	if len(envPaths) > 0 {
		fullEnv := c.FullEnv()
		pathStr := sysutil.ToEnvPATH(envPaths)

		// parse env PATH value
		if len(vars) > 0 && strutil.ContainsByte(pathStr, '$') {
			svRender.SetGetter(func(name string) (val string, ok bool) {
				if val, ok = fullEnv[name]; ok {
					return val, true
				}
				return name, false
			})
			pathStr = svRender.Render(pathStr, vars)
		}

		// 自定义PATH放在最前面, 会优先搜索
		envMap["PATH"] = pathStr + string(filepath.ListSeparator) + sysutil.Getenv("PATH")
	}

	return envMap
}
