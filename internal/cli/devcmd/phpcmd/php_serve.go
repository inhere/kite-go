package phpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// PhpServeCmd instance
var PhpServeCmd = &gcli.Command{
	Name: "serve",
	Desc: "start an php development server",
	Func: func(c *gcli.Command, args []string) error {
		cmd := cmdr.NewCmd("php", "-v")

		return cmd.Run()
	},
}
