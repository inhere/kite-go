package httpcmd

import (
	"github.com/gookit/gcli/v3"
)

// HttpCmd command
var HttpCmd = &gcli.Command{
	Name: "http",
	// Aliases: []string{"h"},
	Desc: "provide some useful tools commands",
	Subs: []*gcli.Command{
		HttpServeCmd,
		SendRequestCmd,
		SendTemplateCmd,
		NewEchoServerCmd(),
		NewFileServerCmd(),
	},
	Config: func(c *gcli.Command) {

	},
}
