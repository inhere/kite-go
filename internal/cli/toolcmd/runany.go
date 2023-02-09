package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/lcproxy"
)

var runOpts = struct {
	listAll bool

	alias, system, search, proxy bool
}{}

// RunAnyCmd instance
var RunAnyCmd = &gcli.Command{
	Name:    "run",
	Desc:    "run custom script command in the `scripts`",
	Aliases: []string{"exec", "script"},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&runOpts.listAll, "list, l", "List information for all scripts or one script")
		c.BoolOpt2(&runOpts.search, "search, s", "Display all matched scripts by the input name")
		c.BoolOpt2(&runOpts.alias, "alias", "dont check and direct run alias command on kite")
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

func runAnything(c *gcli.Command, args []string) error {
	// set proxy ENV
	if runOpts.proxy {
		dump.P(app.App().Lcp)
		app.App().Lcp.Apply(func() {
			c.Infoln("TIP: enabled to set proxy ENV vars, by", envutil.Getenv(lcproxy.HttpKey))
		})
	}

	name := c.Arg("command").String()
	if strutil.IsBlank(name) {
		return c.NewErr("please input a command for run")
	}

	if runOpts.system {
		return cmdr.NewCmd(name, args...).FlushRun()
	}

	return errorx.Raw("TODO")
}
