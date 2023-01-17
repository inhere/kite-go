package comtool

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/rux/handlers"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/web"
	"github.com/inhere/kite/pkg/httpserve"
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

// HttpServe Command
var HttpServe = &gcli.Command{
	Name:    "serve",
	Desc:    "start an http application serve",
	Aliases: []string{"server", "http-serve"},
	Config: func(c *gcli.Command) {
		// bind options
		c.StrOpt(&httpServeOpts.env, "env", "", app.EnvDev, "the application env name")
		c.BoolOpt(&httpServeOpts.debug, "debug", "", true, "the debug mode for run serve")
		c.StrOpt(&httpServeOpts.runtime, "runtime", "", "", "the runtime directory path")

		c.StrVar(&httpServeOpts.host, &gcli.FlagMeta{
			Name:   "host",
			Shorts: []string{"h"},
			Desc:   "host for the start http serve",
			DefVal: "127.0.0.1",
		})
		c.IntVar(&httpServeOpts.port, &gcli.FlagMeta{
			Name:   "port",
			Shorts: []string{"p"},
			Desc:   "port for the start http serve",
			DefVal: 8080,
		})
	},
	Func: func(c *gcli.Command, args []string) error {

		dump.P(httpServeOpts)

		s := httpserve.New()

		r := s.Rux()
		r.Use(handlers.PanicsHandler())

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
