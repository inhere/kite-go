package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/web"
	"github.com/inhere/kite-go/pkg/httpserve"
)

// options for the HttpServe
var httpServeOpts = struct {
	env        string
	host       string
	port       int
	debug      bool
	runtime    string
	staticPath string
}{}

// HttpServeCmd Command
var HttpServeCmd = &gcli.Command{
	Name:    "serve",
	Desc:    "start an http application serve",
	Aliases: []string{"server", "http-serve"},
	Config: func(c *gcli.Command) {
		// bind options
		c.StrOpt(&httpServeOpts.env, "env", "", app.EnvDev, "the application env name")
		c.BoolOpt(&httpServeOpts.debug, "debug", "", true, "the debug mode for run serve")
		c.StrOpt(&httpServeOpts.runtime, "runtime", "", "", "the runtime directory path")

		c.StrVar(&httpServeOpts.host, &gcli.CliOpt{
			Name:   "host",
			Shorts: []string{"h"},
			Desc:   "host for the start http serve",
			DefVal: "127.0.0.1",
		})
		c.IntVar(&httpServeOpts.port, &gcli.CliOpt{
			Name:   "port",
			Shorts: []string{"p"},
			Desc:   "port for the start http serve",
			DefVal: 18080,
		})
	},
	Func: func(c *gcli.Command, args []string) error {
		// dump.P(httpServeOpts)
		s := httpserve.New(httpServeOpts.debug)

		web.AddRoutes(s.Rux())

		// quick start
		// r.Listen("127.0.0.1:18080")
		// apply global pre-handlers
		// http.ListenAndServe(":18080", handlers.HTTPMethodOverrideHandler(r))
		s.Start()
		return nil
	},
}
