package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/kiteext"
	"github.com/inhere/kite/pkg/lcproxy"
)

var runOpts = struct {
	cmdbiz.CommonOpts
	wrapType gflag.EnumString

	listAll, showInfo, search, proxy bool

	alias, plugin, script, system bool
}{}

// RunAnyCmd instance
var RunAnyCmd = &gcli.Command{
	Name:    "run",
	Desc:    "Run any aliases and scripts, as well as plug-ins and system commands",
	Aliases: []string{"exec"},
	Config: func(c *gcli.Command) {
		runOpts.BindCommonFlags(c)
		runOpts.wrapType.SetEnum(kiteext.AllowTypes)

		c.BoolOpt2(&runOpts.listAll, "list, l", "List information for all scripts or one script")
		c.BoolOpt2(&runOpts.showInfo, "show, info, i", "Show information for input alias/script/plugin name")
		c.BoolOpt2(&runOpts.search, "search, s", "Display all matched scripts by the input name")
		c.BoolOpt2(&runOpts.plugin, "plugin", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.alias, "alias", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.script, "script", "dont check and direct run user script on kite")
		c.BoolOpt2(&runOpts.system, "system, sys", "dont check and direct run command on system")
		c.VarOpt(&runOpts.wrapType, "type", "", "wrap shell type for run input script, allow: "+runOpts.wrapType.EnumString())

		c.AddArg("command", "The command for execute, can be with custom arguments")
	},
	Func: runAnything,
	Help: `
## System command

$ kite run ls -al

## Custom scripts

> default in the scripts.yml or dir: $base/scripts

Can use '$@' '$?' at script line. will auto replace to input arguments
examples:

  # scripts.yml
  st: git status
  co: git checkout $@
  br: git branch $?
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

	// set proxy ENV
	if runOpts.proxy {
		app.App().Lcp.Apply(func(lp *lcproxy.LocalProxy) {
			c.Infoln("TIP: enabled to set proxy ENV vars, by", lcproxy.HttpKey, lcproxy.HttpsKey)
			dump.NoLoc(lp)
		})
	}

	// direct run system command
	if runOpts.system {
		c.Infof("TIP: will direct run system command %q (by --system)\n", name)
		return cmdr.NewCmd(name, args...).WorkDirOnNE(runOpts.Workdir).FlushRun()
	}

	// direct run as cmd-alias
	if runOpts.alias {
		c.Infof("TIP: will direct run app command alias %q (by --alias)\n", name)
		return cmdbiz.RunKiteCmdByAlias(name, args)
	}

	ctx := &kiteext.RunCtx{
		Workdir: runOpts.Workdir,
		DryRun:  runOpts.DryRun,
		Type:    runOpts.wrapType.String(),
	}

	// direct run as script
	if runOpts.script {
		c.Infof("TIP: will direct run %q as script name (by --script)\n", name)
		ctx.BeforeFn = func(si *kiteext.ScriptItem) {
			// cliutil.Infof("TIP: %q is a script name, will run it with %v\n", name, args)
			show.AList("Script Context", si)
		}

		return app.Scripts.Run(name, args, ctx)
	}

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

		var si *kiteext.ScriptItem
		si, err = app.Scripts.ScriptItem(name)
		if err != nil {
			return err
		}
		if si != nil {
			show.AList("script info", si)
			return
		}

		si, err = app.Scripts.ScriptFileItem(name)
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
