package kscript

import (
	"fmt"
	"strings"
	"time"

	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/gookit/goutil/x/ccolor"
)

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

	// enhance: if name contains space, as keywords to a search script.
	if strutil.ContainsByte(name, ' ') {
		keywords := strutil.SplitTrimmed(name, " ")
		matched := r.SearchByKeywords(keywords, 2)

		if ln := len(matched); ln == 0 {
			return false, errorx.Rawf("script %q is not exists", name)
		} else if ln > 1 {
			names := maputil.TypedKeys(matched)
			return false, errorx.Rawf("input: %q match more than one script: %s", name, strutil.JoinComma(names))
		}

		// run matched task
		name = maputil.FirstKey(matched)
		if !ctx.Silent {
			ccolor.Greenf("NOTE: match script %q by input keywords %v\n", name, keywords)
		}
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

	// merge env
	envMap := ctx.MergeEnv(r.taskSettings.Env, st.Env)

	// merge env PATH
	if len(r.taskSettings.EnvPaths) > 0 {
		st.EnvPaths = append(st.EnvPaths, r.taskSettings.EnvPaths...)
	}
	if len(st.EnvPaths) > 0 {
		envPaths := append(st.EnvPaths, sysutil.EnvPaths()...)
		envMap["PATH"] = sysutil.ToEnvPATH(envPaths)
	}

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
		"workdir": "",
		"dirname": "",
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

	// workdir
	if strutil.ContainsByte(workdir, '$') {
		workdir = r.renderTaskVars(workdir, vars, ctx)
		vars["workdir"] = workdir
		vars["dirname"] = fsutil.Name(workdir)
	}

	if ctx.Verbose {
		show.AList("Task Vars", vars)
	}

	showIndex := len(st.Cmds) > 1 && !ctx.Silent

	// exec each command
	for idx, line := range st.Cmds {
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

		line = r.renderTaskVars(line, vars, ctx)

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		if showIndex {
			fmt.Printf("-------------------- Run Command #%d --------------------\n", idx+1)
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
func (r *Runner) renderTaskVars(line string, vars map[string]any, ctx *RunCtx) string {
	envs := sysutil.EnvMapWith(ctx.Env)

	rpl.SetGetter(func(name string) (val string, ok bool) {
		// eg: $SHELL
		if r.ParseEnv {
			if val, ok = envs[name]; ok {
				return val, true
			}
		}
		return name, false
	})

	return rpl.Render(line, vars)
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
