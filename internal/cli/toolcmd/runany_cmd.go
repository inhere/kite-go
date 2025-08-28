package toolcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/kscript"
)

type runAnyHandle struct {
	cmdbiz.CommonOpts
	runType  gflag.EnumString // run type
	wrapType gflag.EnumString // shell type
	// input vars
	varMap gflag.KVString
	envMap   gflag.KVString
	// auto find and chdir
	chdir string

	listAll, showInfo, search, verbose bool
}

func (h *runAnyHandle) IsType(name string) bool {
	return h.runType.String() == name
}

var runOpts = runAnyHandle{}

// RunAnyCmd instance
var RunAnyCmd = &gcli.Command{
	Name:    "run",
	Desc: "Run any aliases, script task/file/apps, plugins and system commands",
	Aliases: []string{"exec"},
	Config: func(c *gcli.Command) {
		runOpts.BindCommonFlags2(c)
		runOpts.wrapType.SetEnum(kscript.AllowTypes)
		runOpts.runType.SetEnum([]string{"alias", "script", "ext", "plugin", "system"})

		c.BoolOpt2(&runOpts.listAll, "list, l", "List information for all scripts or one script")
		c.BoolOpt2(&runOpts.showInfo, "show, info, i", "Show information for input alias/script/plugin name")
		c.BoolOpt2(&runOpts.search, "search, s", "Display all matched scripts by the input name")
		c.BoolOpt2(&runOpts.verbose, "verbose, verb", "Display context information on execute")

		c.StrOpt2(&runOpts.chdir, "chdir, cd", "auto find match dir and chdir as workdir")
		c.VarOpt2(&runOpts.envMap, "env,e", "custom set ENV value on run command, format: `KEY=VALUE`")
		c.VarOpt2(&runOpts.varMap, "vars,var", "custom set var value on run command, format: `name=value`")
		c.VarOpt(&runOpts.runType, "type", "t", "direct set type for run input, allow: "+runOpts.runType.EnumString())
		c.VarOpt(&runOpts.wrapType, "shell", "", "wrap shell type for run input script, allow: "+runOpts.wrapType.EnumString())

		c.AddArg("command", "The command for execute, can be with custom arguments")
	},
	Func: runAnything,
	Help: `
## System command

$ kite run ls -al

## Run script task

$ kite run --var key0=value0 --var key1=value1 task_name [args ... for task]

## Custom scripts

> default in the $config/scripts.yml or dir: $base/scripts

Can use '$@' '$*' at script line. will auto replace to input arguments
examples:

  # scripts.yml
  st: git status
  co: git checkout $@
  br: git branch $*
`,
}

func runAnything(c *gcli.Command, args []string) (err error) {
	if runOpts.listAll {
		return listInfos()
	}

	name := c.Arg("command").String()
	if strutil.IsBlank(name) {
		return c.NewErr("please input a command name for run")
	}

	if runOpts.showInfo {
		return showInfo(name)
	}

	wd := runOpts.Workdir
	if wd == "" && runOpts.chdir != "" {
		// cwd := sysutil.Workdir()
		cd, changed := fsutil.SearchNameUpx(sysutil.Workdir(), runOpts.chdir)
		if changed {
			wd = cd
			colorp.Yellowf("TIP: auto find the %q and will chdir to %s\n", runOpts.chdir, cd)
		} else if cd == "" {
			colorp.Warnf("TIP: can not find the %q in %s or parent\n", runOpts.chdir, wd)
		}
	}

	// direct run system command
	if runOpts.IsType("system") {
		colorp.Infof("TIP: will direct run system command %q (by --type=system)\n", name)
		return cmdr.NewCmd(name, args...).WorkDirOnNE(wd).FlushRun()
	}

	// direct run as cmd-alias
	if runOpts.IsType("alias") {
		colorp.Infof("TIP: will direct run app command alias %q (by --type=alias)\n", name)
		return cmdbiz.RunKiteCmdByAlias(name, args)
	}

	ctx := &kscript.RunCtx{
		Workdir: wd,
		Verbose: runOpts.verbose,
		DryRun:  runOpts.DryRun,
		// custom ENV, vars
		Env:  runOpts.envMap.Data(),
		Vars: runOpts.varMap.Data(),
		Type: runOpts.wrapType.String(),
	}

	// direct run as a script
	if runOpts.IsType("script") {
		if runOpts.search {
			ret := app.Scripts.Search(name, args, 10)
			show.AList("Results of search:", ret)
			return nil
		}

		colorp.Infof("TIP: will direct run %q as script name (by --type=script)\n", name)
		ctx.WithNameArgs(name, args)
		cmdbiz.ConfigScriptCtx(ctx)
		return app.Scripts.Run(name, args, ctx)
	}

	// search ...
	if runOpts.search {
		ret := app.Scripts.Search(name, args, 10)
		show.AList("Search scripts:", ret)
		return nil
	}

	// try alias, ext, script, ...
	return cmdbiz.RunAny(name, args, ctx)
}

func showInfo(name string) (err error) {
	if runOpts.IsType("alias") {
		if app.Kas.HasAlias(name) {
			cliutil.Infoln("Alias  :", name)
			cliutil.Infoln("Command:", app.Kas.ResolveAlias(name))
			return
		}
		return errorx.Rawf("app command alias %q is not exists", name)
	}

	if err = app.Scripts.InitLoad(); err != nil {
		return err
	}

	if runOpts.IsType("script") || app.Scripts.IsScriptTask(name) {
		si, err1 := app.Scripts.LoadScriptTaskInfo(name)
		if err1 != nil {
			return err1
		}
		if si != nil {
			show.AList("script task info", si)
			return
		}

		sf, err2 := app.Scripts.LoadScriptFileInfo(name)
		if err2 != nil {
			return err2
		}
		if sf != nil {
			show.AList("script file info", si)
			return
		}
		return errorx.Rawf("input %q is not script or script-file", name)
	}

	return errorx.New("TODO")
}

func listInfos() (err error) {
	if runOpts.IsType("alias") {
		show.AList("command aliases", app.Kas)
		return
	}

	// todo list plugins

	err = app.Scripts.InitLoad()
	if err != nil {
		return err
	}

	if !runOpts.IsType("script") {
		show.AList("command aliases", app.Kas)
		return
	}

	show.AList("loaded script tasks", app.Scripts.RawScriptTasks())
	show.AList("loaded script files", app.Scripts.ScriptFiles())
	return
}

// ScriptCmd command
// var ScriptCmd = &gcli.Command{
// 	Name: "script",
// 	// Aliases: []string{"rand"},
// 	Desc: "list the jump storage data in local",
// 	Config: func(c *gcli.Command) {
// 		// random string(number,alpha,), int(range)
// 	},
// 	Func: func(c *gcli.Command, _ []string) error {
// 		return errorx.New("TODO")
// 	},
// }
