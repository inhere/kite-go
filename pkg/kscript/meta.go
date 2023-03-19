package kscript

import (
	"regexp"
	"sort"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
)

// ScriptInfo struct
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

	//
	// For script file
	//

	// File script file path in Runner.ScriptDirs
	File    string
	BinName string
	FileExt string // eg: .go
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

// InitType on not set
func (si *ScriptInfo) InitType(typ string) {
	if si.Type == "" {
		si.Type = typ
	}
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
