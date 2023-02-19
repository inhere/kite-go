package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cliutil/cmdline"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/lcproxy"
)

var runOpts = struct {
	cmdTyp  string
	listAll bool

	alias, plugin, script, system, search, proxy bool
}{}

// RunAnyCmd instance
var RunAnyCmd = &gcli.Command{
	Name:    "run",
	Desc:    "run custom script command in the `scripts`",
	Aliases: []string{"exec", "script"},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&runOpts.listAll, "list, l", "List information for all scripts or one script")
		c.BoolOpt2(&runOpts.search, "search, s", "Display all matched scripts by the input name")
		c.BoolOpt2(&runOpts.plugin, "plugin", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.alias, "alias", "dont check and direct run alias command on kite")
		c.BoolOpt2(&runOpts.script, "script", "dont check and direct run user script on kite")
		c.BoolOpt2(&runOpts.system, "system, sys", "dont check and direct run command on system")
		c.BoolOpt2(&runOpts.proxy, "proxy,p", "set proxy ENV on run command(config:local_proxy)")

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

var ka maputil.Aliases

func runAnything(c *gcli.Command, args []string) (err error) {
	ka = app.Cfg().StringMap("aliases")

	if runOpts.listAll {
		if runOpts.alias {
			show.AList("command aliases", ka)
			return
		}

		show.AList("command aliases", ka)
		return
	}

	// set proxy ENV
	if runOpts.proxy {
		app.App().Lcp.Apply(func(lp *lcproxy.LocalProxy) {
			c.Infoln("TIP: enabled to set proxy ENV vars, by", lcproxy.HttpKey, lcproxy.HttpsKey)
			dump.NoLoc(lp)
		})
	}

	name := c.Arg("command").String()
	if strutil.IsBlank(name) {
		return c.NewErr("please input a command for run")
	}

	// direct run
	if runOpts.system {
		c.Infof("(by --system) TIP: will direct run system command %q\n", name)
		return cmdr.NewCmd(name, args...).FlushRun()
	}

	// direct run
	if runOpts.alias {
		c.Infof("(by --alias) TIP: will direct run cli command alias %q\n", name)
		return runKiteCmdByAlias(name, args)
	}

	// maybe is kite command alias
	if ka.HasAlias(name) {
		c.Infof("TIP: %q is an cli command alias, will run it with %v\n", name, args)
		return runKiteCmdByAlias(name, args)
	}

	// TODO is script, plugin

	// maybe is system command name
	if sysutil.HasExecutable(name) {
		c.Infof("TIP: %q is a executable file on system, will run it with %v\n", name, args)
		return cmdr.NewCmd(name, args...).FlushRun()
	}
	return errorx.Rawf("%q is not an alias OR script OR plugin OR system command name", name)
}

func runKiteCmdByAlias(name string, inArgs []string) error {
	if !ka.HasAlias(name) {
		return errorx.Newf("kite alias command %q is not found", name)
	}

	str := ka.ResolveAlias(name)
	lp := cmdline.NewParser(str)

	cmd, args := lp.BinAndArgs()
	if len(inArgs) > 0 {
		args = append(args, inArgs...)
	}

	if !app.Cli().HasCommand(cmd) {
		return errorx.Rawf("cli command %q not exist, config in alias: %s", cmd, name)
	}
	return app.Cli().RunCmd(cmd, args)
}
