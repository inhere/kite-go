package cmd

import (
	"github.com/gookit/gcli/v2"
	"github.com/gookit/kite/app"
)

// options for the HttpServe
var httpServeOpts = struct {
	env     string
	runtime string
}{}

var HttpServe = &gcli.Command{
	Name:   "serve",
	UseFor: "start an http application serve",
	Aliases: []string{"server", "http:serve"},
	Config: func(c *gcli.Command) {
		// bind options
		c.StrOpt(&httpServeOpts.env, "env", "", app.EnvDev, "the application env name")
		c.StrOpt(&httpServeOpts.runtime, "runtime", "", "", "the runtime directory path")
	},
}
