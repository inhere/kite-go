package kscript

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/expr-lang/expr"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
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
	//  - usage in a task: ${vars.key}
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
		envMap := maputil.AnyToStrMap(envData)
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
	Desc  string
	Usage string // usage message
	Help  string // long help message

	// Type shell wrap for run the task. allow: sh, bash, zsh
	Type string
	// Alias names for the script task
	Alias []string
	// Silent mode, dont print exec command line and output.
	Silent bool `json:"silent"`

	// Output target. default is stdout
	Output string
	// Vars for run script task.
	//  - task配置中访问: $name
	//  - allow dynamic var: "@sh: git log -1" TODO
	Vars map[string]string
	// Deps task name list. 当前任务依赖的任务名称列表
	Deps []string `json:"deps"`

	// Cmds exec commands list.
	Cmds []*TaskCmd // TODO
	// Args for exec task commands.
	Args []string
	// CmdTimeout for run each command, default is 0.
	CmdTimeout time.Duration

	// Platform limit exec. allow: windows, linux, darwin
	Platform []string
	// PlatformSet 当前系统平台的设置，可以覆盖设置 Type, Cmds
	PlatformSet map[string]any

	// Ext enable extensions: proxy, clip
	Ext string
	// If condition check for run command. eg: sh:test -f .env
	// or see github.com/expr-lang/expr
	If string
}

// ScriptInfo one script task.
type ScriptInfo = ScriptTask

func parseScriptTask(name string, info any, fbType string) (*ScriptTask, error) {
	st := &ScriptTask{Name: name}
	st.ScriptType = TypeTask

	var err error
	switch typVal := info.(type) {
	case string: // one command
		err = st.loadCmdByString(typVal)
	case []string: // as commands
		err = st.loadCmdsByStrings(typVal)
	case []any: // as commands, but support adv command
		err = st.loadCmdsByAnySlice(typVal)
	case map[string]any: // as structured
		err = st.LoadFromMap(typVal)
	default:
		return nil, errorx.Rawf("invalid info of the script task %q, info: %v", name, info)
	}

	st.WithFallbackType(fbType)
	return st, err
}

var (
	// runKeys = []string{"run", "cmds", "cmd", "command", "commands"}
	runKeys = []string{"run", "cmds", "cmd"}
)

func (st *ScriptTask) LoadFromMap(mp map[string]any) error {
	data := maputil.Data(mp)
	st.Type = data.Str("type")
	st.Workdir = data.StrOne("dir", "workdir")
	st.Desc = data.StrOne("desc", "description")

	taskTimeout := data.Str("timeout")
	taskDur, err := timex.ToDuration(taskTimeout)
	if err != nil {
		return errorx.Ef("invalid timeout of the task %q, timeout=%s", st.Name, taskTimeout)
	}
	st.Timeout = taskDur

	cmdTimeout := data.Str("cmd_timeout")
	cmdDur, err := timex.ToDuration(cmdTimeout)
	if err != nil {
		return errorx.Ef("invalid cmd timeout of the task %q, cmd_timeout=%s", st.Name, cmdTimeout)
	}
	st.Timeout = cmdDur

	err = st.loadArgsDefine(data.Get("args"))
	if err != nil {
		return err
	}

	// st.Vars 支持动态变量
	st.Vars = data.StringMap("vars")
	st.Deps = data.StringsOne("deps", "depends")
	// st.Cmds = data.StringsOne("run", "cmd", "cmds")
	cmds := data.One(runKeys...)

	// append set env
	st.Env = data.StringMap("env")
	st.EnvPaths = data.StringsOne("env_path", "env_paths")

	// windows, linux, darwin 每个平台都可以覆盖前面的配置 type, cmds
	subData := data.Sub(runtime.GOOS)
	if len(subData) > 0 {
		st.PlatformSet = subData
		if typStr := subData.Str("type"); typStr != "" {
			st.Type = typStr
		}
		if cmdsVal := subData.One(runKeys...); cmdsVal != nil {
			cmds = cmdsVal
		}
	}

	return st.loadTaskCmdsByAny(cmds)
}

func (st *ScriptTask) loadTaskCmdsByAny(val any) error {
	switch typVal := val.(type) {
	case string: // one command
		return st.loadCmdByString(typVal)
	case []string: // as commands
		return st.loadCmdsByStrings(typVal)
	case []any: // as commands
		return st.loadCmdsByAnySlice(typVal)
	default:
		return errorx.Rawf("invalid cmd info of the task %q, info: (%T)%v", st.Name, typVal, val)
	}
}

func (st *ScriptTask) loadCmdByString(val string) error {
	st.Cmds = append(st.Cmds, newTaskCmd(st, val))
	return nil
}

func (st *ScriptTask) loadCmdsByStrings(ss []string) error {
	for i, cmd := range ss {
		st.Cmds = append(st.Cmds, newTaskCmd2(st, cmd, i))
	}
	return nil
}

func (st *ScriptTask) loadCmdsByAnySlice(anySlice []any) error {
	// as commands
	for i, cmd := range anySlice {
		switch cmdVal := cmd.(type) {
		case string:
			st.Cmds = append(st.Cmds, newTaskCmd2(st, cmdVal, i))
		case map[string]any: // as TaskCmd, 可以设置workdir, run 等
			if len(cmdVal) == 0 {
				return nil
			}

			tc := &TaskCmd{st: st, index: i}
			if err := tc.loadFromMap(cmdVal); err != nil {
				return err
			}

			st.Cmds = append(st.Cmds, tc)
		}
	}

	return nil
}

var argReg = regexp.MustCompile(`\$\d{1,2}`)

// ParseArgs on commands
func (st *ScriptTask) ParseArgs() (args []string) {
	if len(st.Cmds) == 0 {
		return
	}

	// 检测命令是否需要类似shell的参数 eg: echo $1
	str := st.CmdsToString()
	ss := arrutil.Unique(argReg.FindAllString(str, -1))

	sort.Strings(ss)
	return ss
}

// CmdsToString build.
func (st *ScriptTask) CmdsToString(sep ...string) string {
	ln := len(st.Cmds)
	if ln == 0 {
		return ""
	}
	if ln == 1 {
		return st.Cmds[0].Run
	}

	ss := make([]string, 0, ln)
	sepStr := arrutil.FirstOr(sep, ",")

	for _, cmd := range st.Cmds {
		ss = append(ss, cmd.Run)
	}
	return strings.Join(ss, sepStr)
}

func (st *ScriptTask) resolveIfExpr(vars map[string]any) (ok bool) {
	program, err := expr.Compile(st.If, expr.Env(vars))
	if err != nil {
		panic(err) // TODO
	}

	output, err := expr.Run(program, vars)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
	return true
}

func (st *ScriptTask) resolveDynVars(vars map[string]string) (smp map[string]string, err error) {
	if len(vars) == 0 {
		return
	}

	smp = make(map[string]string, len(vars))

	// 解析动态变量值
	for key, val := range vars {
		// eg: @sh:test -v app.go
		if pos := strings.Index(val, ":"); pos > 1 && val[0] == '@' {
			typ := val[1:pos]
			line := strings.TrimSpace(val[pos+1:])
			slog.Debugf("task %s resolveDynVars: %s type=%s, expr=%q", st.Name, key, typ, line)

			switch typ {
			case "sh", "bash", "zsh", "cmd", "pwsh":
				str, err1 := sysutil.ShellExec(line, typ)
				if err1 != nil {
					return nil, err1
				}
				smp[key] = strings.TrimSpace(str)
			case "exec":
				str, err1 := sysutil.ExecLine(line)
				if err1 != nil {
					return nil, err1
				}
				smp[key] = strings.TrimSpace(str)
			default: // as normal string.
				smp[key] = val
			}
			continue
		}

		// normal const value
		smp[key] = val
	}

	return
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

// TaskCmd of the task TODO
type TaskCmd struct {
	st *ScriptTask
	// is reference another task. eg: @task:another_task
	isRef bool
	index int

	// Workdir for run command
	Workdir string
	// Vars for run cmd.
	//  - task配置中访问: $name
	//  - allow dynamic var: "@sh: git log -1" TODO
	Vars map[string]string
	// Env append ENV setting for run
	Env map[string]string
	// Run command line expr for run. eg: go run main.go
	Run string
	// Task refer task name.
	Task string
	// Type wrap for run. Allow: sh, bash, zsh
	Type string
	// If condition expr for run, return true or false
	If string
	// FailMsg custom message on run fail
	FailMsg string
	// Silent mode, dont print exec command line.
	Silent bool `json:"silent"`
	// Timeout for run the command, default is 0.
	Timeout time.Duration
}

func newTaskCmd(st *ScriptTask, run string) *TaskCmd {
	tc := &TaskCmd{st: st}
	tc.loadRun(run)
	return tc
}

func newTaskCmd2(st *ScriptTask, run string, index int) *TaskCmd {
	tc := &TaskCmd{st: st, index: index}
	tc.loadRun(run)
	return tc
}

func (tc *TaskCmd) loadFromMap(mp map[string]any) error {
	data := maputil.Data(mp)

	tc.Type = data.Str("type")
	tc.Task = data.Str("task")
	tc.Vars = data.StringMap("vars")
	tc.Env = data.StringMap("env")
	// more setting
	tc.Silent = data.Bool("silent")
	tc.FailMsg = data.Str("fail_msg")
	tc.Workdir = data.StrOne("workdir", "dir")

	cmdTimeout := data.Str("timeout")
	cmdDur, err := timex.ToDuration(cmdTimeout)
	if err != nil {
		return errorx.Ef("invalid timeout of the task %q command#%d, timeout=%s", tc.st.Name, tc.index, cmdTimeout)
	}

	tc.Timeout = cmdDur
	tc.loadRun(data.StrOne(runKeys...))
	return nil
}

func (tc *TaskCmd) loadRun(run string) {
	if strings.HasPrefix(run, "@task:") {
		tc.isRef = true
		tc.Run = run[6:]
		tc.Task = tc.Run
		return
	}

	if tc.Task != "" {
		tc.isRef = true
		tc.Run = tc.Task
		return
	}

	// first is @ for silent exec
	if strings.HasPrefix(run, "@") {
		tc.Run = run[1:]
		tc.Silent = true
		return
	}

	tc.Run = run
}

func (tc *TaskCmd) appendVars(vars map[string]any) error {
	if len(tc.Vars) == 0 {
		return nil
	}

	// parse dynamic vars
	tcVars, err := tc.st.resolveDynVars(tc.Vars)
	if err != nil {
		return errorx.Rf("task command#%d: %v", tc.index, err)
	}

	// set vars
	for k, v := range tcVars {
		vars[k] = v
	}
	return nil
}
