package kscript

import (
	"fmt"
	"time"

	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/gookit/goutil/timex"
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
	if err = r.InitLoad(); err != nil {
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
	ctx.Args = args

	// ------ try check is task and run it ------
	si, err := r.LoadScriptTaskInfo(name)
	if err != nil {
		return found, err
	}
	if si != nil {
		ccolor.Magentaf("Run Script Task %q Args: %v\n", name, args)
		return found, r.runScriptTask(si, args, ctx)
	}

	// ------ try check is file and run it ------
	sf, err := r.LoadScriptFileInfo(name)
	if err != nil {
		return found, err
	}

	if sf != nil {
		ccolor.Magentaf("Run Script File %q Args: %v\n", name, args)
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

	cmdLn := len(st.Cmds)
	if cmdLn == 0 {
		return errorx.Rawf("empty cmd config for script task %q", ctx.Name)
	}

	needArgs := st.ParseArgs()
	if nln := len(needArgs); len(ctx.Args) < nln {
		ccolor.Println("<mga>Script task contents:</>\n ", st.CmdsToString("\n  "))
		return errorx.Rawf("missing required args for run task %q(need %d)", ctx.Name, nln)
	}

	shell := strutil.OrElse(ctx.Type, st.Type)
	workdir := strutil.OrElse(ctx.Workdir, st.Workdir)

	// build context vars 用于渲染 ENV, PATH, task-run 等的变量
	vars, err := r.buildTaskRenderVars(st, ctx)
	if err != nil {
		return err
	}
	vars["cur_dir"] = sysutil.Workdir()

	// merge ENV and PATH
	ctx.MergeEnv(r.taskSettings.Env, st.Env)
	if len(r.taskSettings.EnvPaths) > 0 {
		st.EnvPaths = append(st.EnvPaths, r.taskSettings.EnvPaths...)
	}

	// parse vars in ENV
	envMap := ctx.ParseVarInEnv(st.EnvPaths, vars)

	// parse workdir
	if strutil.ContainsByte(workdir, '$') {
		workdir = r.renderTaskVars(workdir, vars, ctx)
		vars["workdir"] = workdir
		vars["dirname"] = fsutil.Name(workdir)
	}

	// append args to vars
	ctx.AppendArgsToVars(vars)
	ccolor.Magentaln("CURRENT DIR:", sysutil.Workdir())
	if ctx.Verbose {
		show.AList("Task Vars", vars)
	}

	// 先执行 deps 任务
	if len(st.Deps) > 0 {
		for _, depTask := range st.Deps {
			ccolor.Magentaln("Run Depends Task:", depTask)

			dst, err1 := r.LoadScriptTaskInfo(depTask)
			if err1 != nil {
				return errorx.Rf("task %s: load dep task %q info fail: %v", st.Name, depTask, err)
			}
			if dst == nil {
				return errorx.Rawf("task %s: the dep task %q not found", st.Name, depTask)
			}

			if err2 := r.runScriptTask(dst, inArgs, ctx); err2 != nil {
				return err2
			}
		}
	}

	startTime := timex.Now().T()
	showIndex := cmdLn > 1 && !ctx.Silent

	// exec each command
	for idx, tc := range st.Cmds {
		if len(tc.Run) == 0 {
			continue
		}

		// redirect runs another task
		if tc.isRef {
			name := tc.Run
			ccolor.Magentaln("Run Refer Task:", name)
			osi, err1 := r.LoadScriptTaskInfo(name)
			if err1 != nil {
				return err1
			}
			if osi == nil {
				return errorx.Rawf("task %q: reference script task %q not found", st.Name, name)
			}

			if err2 := r.runScriptTask(osi, inArgs, ctx); err2 != nil {
				return err2
			}
			continue
		}

		// 加载 command 独有的变量
		if err1 := tc.appendVars(vars); err1 != nil {
			return err1
		}

		line := r.renderTaskVars(tc.Run, vars, ctx)
		// workdir for cmd
		cmdDir := strutil.OrElse(tc.Workdir, workdir)
		if strutil.ContainsByte(cmdDir, '$') {
			cmdDir = r.renderTaskVars(cmdDir, vars, ctx)
			vars["workdir"] = cmdDir
			vars["dirname"] = fsutil.Name(cmdDir)
		}

		var cmd *cmdr.Cmd
		if shell != "" {
			cmd = cmdr.NewCmd(shell, "-c", line)
		} else {
			cmd = cmdr.NewCmdline(line)
		}

		if showIndex {
			fmt.Printf("--------------------------- task command #%d ---------------------------\n", idx+1)
		}
		err2 := cmd.WorkDirOnNE(cmdDir).WithDryRun(ctx.DryRun).AppendEnv(envMap).PrintCmdline2().FlushRun()
		if err2 != nil {
			return err2
		}
	}

	ccolor.Infof(" ✅  Task %s: all task commands done. Take time: %s\n", st.Name, timex.Now().Diff(startTime))
	return nil
}

func (r *Runner) buildTaskRenderVars(st *ScriptTask, ctx *RunCtx) (map[string]any, error) {
	// build context vars
	data := map[string]any{
		"workdir": "",
		"dirname": "",
	}

	// st.Vars 需要支持动态变量
	stVars, err := st.resolveDynVars(st.Vars)
	if err != nil {
		return nil, err
	}

	// 当前task配置的和ctx输入变量，放在顶级直接访问
	topVars := maputil.MergeStrMap(ctx.Vars, stVars)
	for k, v := range topVars {
		data[k] = v
	}
	// dump.P(ctx, topVars)

	// 内置扩展变量
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

	if ctx.AppendVarsFn != nil {
		data = ctx.AppendVarsFn(data)
	}
	return data, nil
}

// 使用简单的模板渲染，支持链式语法变量替换，环境变量，默认值等 - 无法同时支持 $var ${var_name}
// var rpl = textutil.NewVarReplacer("$").WithParseEnv().WithParseDefault()
// 专门实现的类似 php, shell 的字符串表达式处理
var rpl = textutil.NewStrVarRenderer()

// process vars and env
func (r *Runner) renderTaskVars(line string, vars map[string]any, ctx *RunCtx) string {
	envs := ctx.FullEnv()

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
