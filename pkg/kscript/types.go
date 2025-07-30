package kscript

import (
	"time"

	"github.com/gookit/goutil/maputil"
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
	Env map[string]string
	// EnvPaths custom prepend set ENV PATH.
	EnvPaths []string
	// Timeout for run a script, default is 0.
	Timeout time.Duration
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
	Env map[string]string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`

	// BeforeFn hook. si: ScriptTask | ScriptApp | ScriptFile
	BeforeFn func(si any, ctx *RunCtx)
	// AppendVarsFn hook for run task. eg: gvs, paths, kite
	AppendVarsFn func(data map[string]any) map[string]any
}

// EnsureCtx to
func EnsureCtx(ctx *RunCtx) *RunCtx {
	if ctx == nil {
		return &RunCtx{}
	}
	return ctx
}

// WithName to ctx
func (c *RunCtx) WithName(name string) *RunCtx {
	c.Name = name
	return c
}

// MergeEnv and returns
func (c *RunCtx) MergeEnv(mps ...map[string]string) map[string]string {
	for _, mp := range mps {
		c.Env = maputil.MergeStrMap(mp, c.Env)
	}
	return c.Env
}
