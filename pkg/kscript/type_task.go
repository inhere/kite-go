package kscript

import (
	"regexp"
	"sort"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

//
// endregion
// region T: TaskSettings
//

// Variable dynamic variable definition. 动态变量 TODO
type Variable struct {
	// Type of variable, allow: sh, bash, zsh, go or empty.
	Type  string
	Expr string
	Value string
}

// TaskSettings 可以通过 script task 文件中的 "__settings" 调整设置
type TaskSettings struct {
	// Env append set ENV for all script tasks
	Env map[string]string `json:"env"`
	// EnvPaths prepend set ENV PATHs for all script tasks
	EnvPaths []string `json:"env_paths"`

	// DefaultGroup default group name for use Groups. will merge default group data to Vars
	DefaultGroup string `json:"default_group"`
	// Vars built in vars map. group name: vars
	//  - usage: ${vars.key}
	Vars map[string]string `json:"vars"`
	// Grouped vars map.
	//  - group name => map[string]string grouped var map.
	Groups comdef.L2StrMap `json:"groups"`
}

func (ts *TaskSettings) loadData(data map[string]any) {
	// dump.P("task settings:", data)

	if varsData, ok1 := data["vars"]; ok1 {
		if varsMap, ok2 := varsData.(map[string]any); ok2 {
			strMap := maputil.ToStringMap(varsMap)
			ts.Vars = maputil.MergeStrMap(strMap, ts.Vars)
		}
	}

	if groupsData, ok1 := data["groups"]; ok1 {
		if groupsMap, ok2 := groupsData.(map[string]any); ok2 {
			l2StrMap := maputil.ToL2StrMap(groupsMap)
			ts.Groups = maputil.MergeL2StrMap(ts.Groups, l2StrMap)
		}
	}

	if defGroupVal, ok := data["default_group"]; ok {
		defGroup := strutil.SafeString(defGroupVal)
		ts.DefaultGroup = defGroup

		if defGroupData, ok := ts.Groups[defGroup]; ok {
			ts.Vars = maputil.MergeStrMap(defGroupData, ts.Vars)
		}
	}

	// env
	if envData, ok := data["env"]; ok {
		envMap := maputil.TryStringMap(envData)
		if envMap != nil {
			ts.Env = maputil.MergeStrMap(ts.Env, envMap)
		}
	}

	// env PATH
	if value, ok := data["env_path"]; ok {
		if envPaths := tryToStrings(value); envPaths != nil {
			ts.EnvPaths = append(ts.EnvPaths, envPaths...)
		}
	}
	if value, ok := data["env_paths"]; ok {
		if envPaths := tryToStrings(value); envPaths != nil {
			ts.EnvPaths = append(ts.EnvPaths, envPaths...)
		}
	}
}

func tryToStrings(value any) []string {
	if str, ok := value.(string); ok {
		return strutil.SplitTrimmed(str, ",")
	}
	return arrutil.AnyToStrings(value)
}

//
// endregion
// region T: ScriptTask
//

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
	// Vars for run script.
	//  - allow exec a command line TODO
	Vars map[string]string
	// Ext enable extensions: proxy, clip
	Ext string
	// Deps commands list. 当前任务依赖的任务名称列表
	Deps []string `json:"deps"`

	// Cmds exec commands list.
	Cmds []string
	// Cmds []*Command
	// Args for exec task commands.
	Args []string

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

// ScriptInfo one script task.
type ScriptInfo = ScriptTask

func parseScriptTask(name string, info any, fbType string) (*ScriptTask, error) {
	st := &ScriptTask{Name: name}
	st.ScriptType = TypeTask

	switch typVal := info.(type) {
	case string: // one command
		st.Cmds = []string{typVal}
	case []string: // as commands
		st.Cmds = typVal
	case []any: // as commands
		st.Cmds = arrutil.SliceToStrings(typVal)
	case map[string]any: // as structured
		data := maputil.Data(typVal)
		st.Type = data.Str("type")
		st.Workdir = data.StrOne("dir", "workdir")
		st.Desc = data.StrOne("desc", "description")

		err := st.loadArgsDefine(data.Get("args"))
		if err != nil {
			return nil, err
		}

		st.Vars = data.StringMap("vars")
		st.Deps = data.StringsOne("deps", "depends")
		st.Cmds = data.StringsOne("run", "cmd", "cmds")

		// append set env
		st.Env = data.StringMap("env")
		st.EnvPaths = data.StringsOne("env_path", "env_paths")

		// TODO override by os platform
		// osName := runtime.GOOS // windows, linux, darwin
		// data.SubMap(osName) // 每个平台都可以覆盖前面的配置
	default:
		return nil, errorx.Rawf("invalid info of the script task %q, info: %v", name, info)
	}

	st.WithFallbackType(fbType)
	return st, nil
}

func (st *ScriptTask) LoadFromMap(data map[string]any) error {
	return nil
}

var argReg = regexp.MustCompile(`\$\d{1,2}`)

// ParseArgs on commands
func (st *ScriptTask) ParseArgs() (args []string) {
	if len(st.Cmds) == 0 {
		return
	}

	// 检测命令是否需要类似shell的参数 eg: echo $1
	str := strings.Join(st.Cmds, ",")
	ss := arrutil.Unique(argReg.FindAllString(str, -1))

	sort.Strings(ss)
	return ss
}

// WithFallbackType on not setting.
func (st *ScriptTask) WithFallbackType(typ string) *ScriptTask {
	if st.Type == "" {
		st.Type = typ
	}
	return st
}

// args type: string, strings
func (st *ScriptTask) loadArgsDefine(args any) error {
	if args == nil {
		return nil
	}

	switch typVal := args.(type) {
	case string: // desc
		st.Args = []string{typVal}
	case []string: // desc list
		st.Args = typVal
	case []any: // desc list
		st.Args = arrutil.SliceToStrings(typVal)
	// case map[string]string: // name with desc. ERROR: map cannot be ordered
	// 	st.Args = typVal
	default:
		return errorx.Rawf("invalid args config for %q", st.Name)
	}
	return nil
}

//
// endregion
// region T: task command
//

// Command of the task TODO
type Command struct {
	st *ScriptTask
	// is reference another command. eg: @task:another_task
	isRef bool

	// Workdir for run command
	Workdir string
	// Vars for run cmd. allow exec a command line TODO
	Vars map[string]string
	// Env setting for run
	Env map[string]string
	// Run command line expr for run. eg: go run main.go
	Run string
	// Type wrap for run. Allow: sh, bash, zsh
	Type string
	// Msg on run fail
	Msg string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
}
