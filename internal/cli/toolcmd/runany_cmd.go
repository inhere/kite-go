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
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/kscript"
)

var runOpts = struct {
	cmdbiz.CommonOpts
	wrapType gflag.EnumString
	envMap   gflag.KVString
	chdir    string // auto find and chdir

	listAll, showInfo, search, verbose bool

	alias, plugin, script, system bool
}{}

// RunAnyCmd instance
var RunAnyCmd = &gcli.Command{
	Name:    "run",
	Desc:    "Run any aliases and scripts, as well as plug-ins and system commands",
	Aliases: []string{"exec"},
	Config: func(c *gcli.Command) {
		runOpts.BindCommonFlags(c)
		runOpts.wrapType.SetEnum(kscript.AllowTypes)

		c.BoolOpt2(&runOpts.listAll, "list, l", "List information for all scripts or one script")
		c.BoolOpt2(&runOpts.showInfo, "show, info, i", "Show information for input alias/script/plugin name")
		c.BoolOpt2(&runOpts.search, "search, s", "Display all matched scripts by the input name")
		c.BoolOpt2(&runOpts.verbose, "verbose, v", "Display context information on execute")

		c.BoolOpt2(&runOpts.plugin, "plugin", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.alias, "alias", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.script, "script", "dont check and direct run user script on kite")
		c.BoolOpt2(&runOpts.system, "system, sys", "dont check and direct run command on system")
		c.StrOpt2(&runOpts.chdir, "chdir, cd", "auto find match dir and chdir as workdir")

		c.VarOpt2(&runOpts.envMap, "env,e", "custom set ENV value on run command, format: `KEY=VALUE`")
		c.VarOpt(&runOpts.wrapType, "type", "", "wrap shell type for run input script, allow: "+runOpts.wrapType.EnumString())

		c.AddArg("command", "The command for execute, can be with custom arguments")
	},
	Func: runAnything,
	Help: `
## System command

$ kite run ls -al

## Custom scripts

> default in the scripts.yml or dir: $base/scripts

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
	if runOpts.chdir != "" {
		cd, changed := fsutil.SearchNameUpx(wd, runOpts.chdir)
		if changed {
			wd = cd
			colorp.Yellowf("TIP: auto find the %q and will chdir to %s", runOpts.chdir, cd)
		} else if cd == "" {
			colorp.Warnf("TIP: can not find the %q in %s or parent", runOpts.chdir, wd)
		}
	}

	if runOpts.verbose {
		colorp.Infoln("(verbose=true) Run command on workdir:", wd)
	}

	// direct run system command
	if runOpts.system {
		colorp.Infof("TIP: will direct run system command %q (by --system)\n", name)
		return cmdr.NewCmd(name, args...).WorkDirOnNE(wd).FlushRun()
	}

	// direct run as cmd-alias
	if runOpts.alias {
		colorp.Infof("TIP: will direct run app command alias %q (by --alias)\n", name)
		return cmdbiz.RunKiteCmdByAlias(name, args)
	}

	ctx := &kscript.RunCtx{
		Workdir: wd,
		Verbose: runOpts.verbose,
		DryRun:  runOpts.DryRun,
		Type:    runOpts.wrapType.String(),
	}

	// direct run as script
	if runOpts.script {
		colorp.Infof("TIP: will direct run %q as script name (by --script)\n", name)

		if runOpts.search {
			ret := app.Scripts.Search(name, args, 10)
			show.AList("Results of search:", ret)
			return nil
		}

		if runOpts.verbose {
			ctx.BeforeFn = func(si *kscript.ScriptInfo, ctx *kscript.RunCtx) {
				// cliutil.Infof("TIP: %q is a script name, will run it with %v\n", name, args)
				show.AList("Script Info", si)
				show.AList("Run Context", ctx)
			}
		}

		return app.Scripts.Run(name, args, ctx)
	}

	// TODO search ...

	// try alias, script, ...
	return cmdbiz.RunAny(name, args, ctx)
}

func showInfo(name string) (err error) {
	if runOpts.alias {
		if cmdbiz.Kas.HasAlias(name) {
			cliutil.Infoln("Alias  :", name)
			cliutil.Infoln("Command:", cmdbiz.Kas.ResolveAlias(name))
			return
		}
		return errorx.Rawf("app command alias %q is not exists", name)
	}

	if runOpts.script {
		if err = app.Scripts.InitLoad(); err != nil {
			return err
		}

		var si *kscript.ScriptInfo
		si, err = app.Scripts.ScriptDefineInfo(name)
		if err != nil {
			return err
		}
		if si != nil {
			show.AList("script info", si)
			return
		}

		si, err = app.Scripts.ScriptFileInfo(name)
		if err != nil {
			return err
		}
		if si != nil {
			show.AList("script file info", si)
			return
		}
		return errorx.Rawf("input %q is not script or script-file", name)
	}

	return errorx.New("TODO")
}

func listInfos() (err error) {
	if runOpts.alias {
		show.AList("command aliases", cmdbiz.Kas)
		return
	}

	// todo list plugins

	if runOpts.script {
		err = app.Scripts.InitLoad()
		if err != nil {
			return err
		}
		// dump.P(app.Scripts)
		show.AList("loaded scripts", app.Scripts.DefinedScripts())
		show.AList("loaded script files", app.Scripts.ScriptFiles())
		return
	}

	err = app.Scripts.InitLoad()
	if err != nil {
		return err
	}

	show.AList("command aliases", cmdbiz.Kas)
	show.AList("loaded scripts", app.Scripts.DefinedScripts())
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
