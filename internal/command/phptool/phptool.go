package phptool

import (
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/cmdutil"
)

var PhpToolsCmd = &gcli.Command{
	Name: "php",
	Desc: "some php tools command",
	Subs: []*gcli.Command{
		PhpInfo,
		PhpServe,
	},
}

var PhpInfo = &gcli.Command{
	Name: "info",
	Desc: "system info for php",
	Func: func(c *gcli.Command, args []string) error {
		cmd := cmdutil.NewCmd()
		cmd.SetBinArgs("php", "-v")

		return cmd.Run()
	},
}
