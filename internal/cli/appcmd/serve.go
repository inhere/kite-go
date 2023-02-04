package appcmd

import "github.com/gookit/gcli/v3"

// BackendServeCmd kite backend background server
var BackendServeCmd = &gcli.Command{
	Name:    "serve",
	Aliases: []string{"be-serve", "server"},
	Desc:    "kite backend serve application",
}
