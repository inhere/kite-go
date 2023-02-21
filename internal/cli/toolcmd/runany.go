package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/biz/cmdbiz"
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

func runAnything(c *gcli.Command, args []string) (err error) {
	if runOpts.listAll {
		if runOpts.alias {
			show.AList("command aliases", cmdbiz.Kas)
			return
		}

		show.AList("command aliases", cmdbiz.Kas)
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
		return cmdbiz.RunKiteCmdByAlias(name, args)
	}

	// try alias, script, ...
	return cmdbiz.RunAny(name, args)
}
