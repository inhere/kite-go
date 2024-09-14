package kscript

import (
	"regexp"
	"sort"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
)

// ScriptInfo struct Info, Meta
type ScriptInfo struct {

	//
	// For define script
	//

	// Type wrap for run script. allow: sh, bash, zsh
	Type string

	// Workdir for run script
	Workdir string
	// Platform limit. allow: windows, linux, darwin
	Platform []string
	// Output target. default is stdout
	Output string
	// Vars for run script. allow exec a command line TODO
	Vars map[string]string
	// Ext enable ext: proxy, clip
	Ext string
	// Deps commands list
	Deps []string

	// Name for script
	Name string
	// Desc message
	Desc string
	// Env setting for run
	Env map[string]string
	// Args script args definition.
	Args, ArgNames []string
	// Cmds commands list
	Cmds []string

	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
	// PreCond for run script. eg: test -f .env
	PreCond string

	//
	// For script file
	//

	// File script file path in Runner.ScriptDirs
	File    string
	BinName string
	FileExt string // eg: .go
}

func newDefinedScriptInfo(name string, info any, fbType string) (*ScriptInfo, error) {
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

// IsFile script
func (si *ScriptInfo) IsFile() bool {
	return si.File != ""
}

// IsDefined script
func (si *ScriptInfo) IsDefined() bool {
	return si.File == ""
}

// WithFallbackType on not set
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

// CmdInfo struct
type CmdInfo struct {
	si *ScriptInfo
	// Workdir for run script
	Workdir string
	// Vars for run cmd. allow exec a command line TODO
	Vars map[string]string
	// Env setting for run
	Env map[string]string
	// Line command line for run. eg: go run main.go
	Line string
	// Type wrap for run. allow: sh, bash, zsh
	Type string
	// Msg on run fail
	Msg string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
}

// Variable struct
type Variable struct {
	// Type of variable, allow: sh, bash, zsh or empty.
	Type  string
	Value string
}

// RunCtx definition
type RunCtx struct {
	// Name for script run
	Name string
	Type string

	// Verbose show more info on run
	Verbose bool
	// DryRun script
	DryRun bool
	// Workdir for run script
	Workdir string
	// Env setting for run
	Env map[string]string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`

	// BeforeFn hook
	BeforeFn func(si *ScriptInfo, ctx *RunCtx)
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
