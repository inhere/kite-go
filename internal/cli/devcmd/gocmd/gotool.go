package gocmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// GoToolsCmd instance
var GoToolsCmd = &gcli.Command{
	Name: "go",
	Desc: "some go tools command",
	Subs: []*gcli.Command{
		AwesomeGoCmd,
		ListBinCmd,
		GoInfoCmd,
	},
}

// GoInfoCmd refer https://github.com/lucor/goinfo
var GoInfoCmd = &gcli.Command{
	Name: "info",
	Desc: "system info for go",
}

// ListBinCmd instance
var ListBinCmd = &gcli.Command{
	Name: "list-bin",
	Desc: "start an php development server",
	Func: func(c *gcli.Command, args []string) error {
		c.Infoln("list bin in $GOPATH/bin")

		path := sysutil.ExpandPath("${GOPATH}/bin")
		cmd := cmdr.NewCmd("ls", path)

		return cmd.Run()
	},
}
