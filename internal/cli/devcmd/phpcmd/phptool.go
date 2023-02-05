package phpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// PhpToolsCmd instance
var PhpToolsCmd = &gcli.Command{
	Name: "php",
	Desc: "some php tools command",
	Subs: []*gcli.Command{
		PhpInfoCmd,
		PhpServeCmd,
	},
}

// PhpInfoCmd instance
var PhpInfoCmd = &gcli.Command{
	Name: "info",
	Desc: "system info for php",
	Func: func(c *gcli.Command, args []string) error {
		cmd := cmdr.NewCmd("php", "-v")

		return cmd.Run()
	},
}
