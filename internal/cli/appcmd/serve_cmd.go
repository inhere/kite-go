package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// BackendServeCmd kite backend background server
var BackendServeCmd = &gcli.Command{
	Name:    "serve",
	Aliases: []string{"be-serve", "server"},
	Desc:    "kite backend serve application",
	Func: func(c *gcli.Command, args []string) error {
		return errorx.New("todo")
	},
}
