package phptool

import (
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/cmdutil"
)

var PhpServe = &gcli.Command{
	Name: "serve",
	Desc: "start an php development server",
	Func: func(c *gcli.Command, args []string) error {
		cmd := cmdutil.NewCmd()
		cmd.SetBinArgs("php", "-v")

		return cmd.Run()
	},
}
