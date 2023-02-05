package mdcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

var MarkdownServeCmd = &gcli.Command{
	Name:    "serve",
	Aliases: []string{"server"},
	Desc:    "convert an markdown table to create DB table SQL",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
