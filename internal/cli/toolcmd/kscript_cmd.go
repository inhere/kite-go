package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/kscript"
)

type ksRunner struct {
	cmdbiz.CommonOpts
	wrapType gflag.EnumString
	envMap   gflag.KVString
	chdir    string // auto find and chdir

	listAll, showInfo, search, verbose bool
}

// NewKScriptCmd create a command instance
func NewKScriptCmd() *gcli.Command {
	var rr = ksRunner{}
	return rr.NewCmd()
}

func (rr ksRunner) NewCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "script",
		Aliases: []string{"run-s", "xs", "rs"},
		Desc:    "Run the kite script task/command by name",
		Config:  rr.Config,
		Func:    rr.Run,
	}
}

func (rr ksRunner) Config(c *gcli.Command) {
	rr.BindCommonFlags(c)
	rr.wrapType.SetEnum(kscript.AllowTypes)

	c.StrOpt2(&rr.chdir, "chdir, cd", "auto find match dir and chdir as workdir")
	c.VarOpt2(&rr.envMap, "env,e", "custom set ENV value on run command, format: `KEY=VALUE`")
	c.VarOpt(&rr.wrapType, "type", "", "wrap shell type for run input script, allow: "+rr.wrapType.EnumString())

	c.BoolOpt2(&rr.listAll, "list, l", "List information for all scripts or one script")
	c.BoolOpt2(&rr.showInfo, "show, info, i", "Show information for input alias/script/plugin name")
	c.BoolOpt2(&rr.search, "search, s", "Display all matched scripts by the input name")
	c.BoolOpt2(&rr.verbose, "verbose, v", "Display context information on execute")

	c.AddArg("command", "The task/command for execute, can be with custom arguments")
}

func (rr *ksRunner) Run(c *gcli.Command, _ []string) error {
	return errorx.New("TODO")
}
