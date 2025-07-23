package kscript

import (
	"regexp"
	"sort"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/errorx"
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
	// Env setting for run the file/app/task
	Env map[string]string
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
// region T: script_task
//

// TaskSettings 可以通过 script task 文件中的 "__settings" 调整设置
type TaskSettings struct {
	// Vars built in vars map. group name: vars
	Vars   map[string]string `json:"vars"`
	Groups comdef.L2StrMap   `json:"groups"`
}

func (ts *TaskSettings) loadData(data map[string]any) {
	if varsData, ok1 := data["vars"]; ok1 {
		if varsMap, ok2 := varsData.(map[string]string); ok2 {
			ts.Vars = maputil.MergeSMap(varsMap, ts.Vars, false)
		}
	}

	if groupsData, ok1 := data["groups"]; ok1 {
		if groupsMap, ok2 := groupsData.(map[string]map[string]string); ok2 {
			ts.Groups = maputil.MergeL2StrMap(ts.Groups, groupsMap)
		}
	}
}

type ScriptTasks struct {
	tasks map[string]*ScriptTask

	// settings for all script tasks.
	//
	// eg:
	//  - vars: map[string]string built in vars map. group name: vars
	//  - group: map[string]map[string]string grouped var map.
	settings TaskSettings
}

// ScriptTask for one script task.
type ScriptTask struct {
	ScriptMeta

	// Name for the script task
	Name string
	// Desc message
	Desc string
	// Type shell wrap for run the script. allow: sh, bash, zsh
	Type string
	// Alias names for the script task
	Alias []string

	// Platform limit. allow: windows, linux, darwin
	Platform []string
	// Output target. default is stdout
	Output string
	// Vars for run script. allow exec a command line TODO
	Vars map[string]string
	// Ext enable extensions: proxy, clip
	Ext string
	// Deps commands list
	Deps []string `json:"deps"`

	// Cmds exec commands list
	Cmds []string
	// Args script task args definition.
	Args, ArgNames []string

	// CmdLinux command lines on different OS. will override the Cmds
	CmdLinux   []string
	CmdDarwin  []string
	CmdWindows []string

	// Silent mode, dont print exec command line and output.
	Silent bool `json:"silent"`
	// IfCond check for run command. eg: sh:test -f .env
	// or see github.com/hashicorp/go-bexpr
	IfCond string
}

// ScriptInfo one script. or TODO ScriptTask, ScriptMeta, ScriptEntry struct
type ScriptInfo struct {

	//
	// For define script
	//

	// Type shell wrap for run the script. allow: sh, bash, zsh
	Type string

	// Workdir for run script, default is current dir.
	Workdir string
	// Platform limit. allow: windows, linux, darwin
	Platform []string
	// Output target. default is stdout
	Output string
	// Vars for run script. allow exec a command line TODO
	Vars map[string]string
	// Ext enable extensions: proxy, clip
	Ext string
	// Deps commands list
	Deps []string

	// Name for the script
	Name string
	// Desc message
	Desc string
	// Env setting for run
	Env map[string]string
	// Args script args definition.
	Args, ArgNames []string
	// Cmds commands list
	Cmds []string

	// Silent mode, dont print exec command line and output.
	Silent bool `json:"silent"`
	// IfCond check for run command. eg: sh:test -f .env
	// or see github.com/hashicorp/go-bexpr
	IfCond string

	//
	// For script file
	//

	// File script file path in Runner.ScriptDirs
	File    string
	BinName string
	FileExt string // eg: .go
}

func parseScriptTask(name string, info any, fbType string) (*ScriptInfo, error) {
	si := &ScriptInfo{Name: name}

	switch typVal := info.(type) {
	case string: // one command
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
		si.Vars = data.StringMap("vars")
		si.Env = data.StringMap("env")
	default:
		return nil, errorx.Rawf("invalid config of the script %q", name)
	}

	si.WithFallbackType(fbType)
	return si, nil
}

var argReg = regexp.MustCompile(`\$\d{1,2}`)

// ParseArgs on commands
func (si *ScriptInfo) ParseArgs() (args []string) {
	if len(si.Cmds) == 0 {
		return
	}

	str := strings.Join(si.Cmds, ",")
	ss := arrutil.Unique(argReg.FindAllString(str, -1))

	sort.Strings(ss)
	return ss
}

// WithFallbackType on not setting.
func (si *ScriptInfo) WithFallbackType(typ string) *ScriptInfo {
	if si.Type == "" {
		si.Type = typ
	}
	return si
}

// args type: string, strings
func (si *ScriptInfo) loadArgsDefine(args any) error {
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

//
// endregion
// region T: var and group
//

// CmdInfo struct TODO
type CmdInfo struct {
	si *ScriptInfo
	// Workdir for run a script
	Workdir string
	// Vars for run cmd. allow exec a command line TODO
	Vars map[string]string
	// Env setting for run
	Env map[string]string
	// Line command line expr for run. eg: go run main.go
	Line string
	// Type wrap for run. Allow: sh, bash, zsh
	Type string
	// Msg on run fail
	Msg string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
}

// Variable dynamic variable definition.
type Variable struct {
	// Type of variable, allow: sh, bash, zsh, go or empty.
	Type  string
	Expr string
	Value string
}

// GroupVars definition
type GroupVars struct {
	// path separator. default: ":"
	pathSep string
	// default group name
	defaultGroup string
	// key is group name.
	data map[string]comdef.StrMap
}

// Get value by key and group
func (gv *GroupVars) Get(group, key string) (string, bool) {
	if gv.data == nil {
		return "", false
	}
	if group == "" {
		group = gv.defaultGroup
	}

	v, ok := gv.data[group][key]
	return v, ok
}

//
// endregion
// region T: run context
//

// RunCtx definition
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
	// Vars for run cmd.
	Vars map[string]string
	// Env setting for run
	Env map[string]string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`

	// BeforeFn hook. si: ScriptTask | ScriptApp | ScriptFile
	BeforeFn func(si any, ctx *RunCtx)
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
func (c *RunCtx) MergeEnv(mp map[string]string) map[string]string {
	return maputil.MergeSMap(mp, c.Env, false)
}
