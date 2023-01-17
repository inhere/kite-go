package devcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// DevToolsCmd command
var DevToolsCmd = &gcli.Command{
	Name:    "dev",
	Aliases: []string{"dt", "devtool"},
	Desc:    "provide some useful dev tools commands",
	Subs: []*gcli.Command{
		HotReloadServe,
	},
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
