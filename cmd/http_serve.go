package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/kite/app"
	"github.com/gookit/kite/pkg/httpserve"
	"github.com/gookit/kite/web"
	"github.com/gookit/rux"
	"github.com/gookit/rux/handlers"
)

// options for the HttpServe
var httpServeOpts = struct {
	env     string
	port    int
	debug    bool
	runtime string
	staticPath string
}{}

// HttpServe Command
var HttpServe = &gcli.Command{
	Name:   "serve",
	UseFor: "start an http application serve",
	Aliases: []string{"server", "http:serve"},
	Config: func(c *gcli.Command) {
		// bind options
		c.StrOpt(&httpServeOpts.env, "env", "", app.EnvDev, "the application env name")
		c.BoolOpt(&httpServeOpts.debug, "debug", "", true, "the debug mode for run serve")
		c.StrOpt(&httpServeOpts.runtime, "runtime", "", "", "the runtime directory path")
	},
	Func: func(c *gcli.Command, args []string) error {
		r := httpserve.NewServe()
		r.Use(handlers.PanicsHandler())

		// handle error
		r.OnError = func(c *rux.Context) {
			if err := c.FirstError(); err != nil {
				color.Error.Println(err)
				c.HTTPError(err.Error(), 400)
				return
			}
		}

		if httpServeOpts.debug {
			r.Use(handlers.RequestLogger())
		}

		web.AddRoutes(r)

		// quick start
		r.Listen("127.0.0.1:18080")
		// apply global pre-handlers
		// http.ListenAndServe(":18080", handlers.HTTPMethodOverrideHandler(r))

		return nil
	},
}
