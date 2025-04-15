package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/rux"
	"github.com/gookit/rux/pkg/handlers"
	"github.com/gookit/rux/pkg/render"
)

// NewEchoServerCmd instance
func NewEchoServerCmd() *gcli.Command {
	var esOpts = struct {
		port uint
	}{}

	return &gcli.Command{
		Name:    "echo-server",
		Desc:    "start an simple echo http server",
		Aliases: []string{"echo-serve", "echo"},
		Config: func(c *gcli.Command) {
			c.UintOpt(&esOpts.port, "port", "P", 0, "custom the echo server port, default will use random `port`")
		},
		Func: func(c *gcli.Command, args []string) error {
			if esOpts.port < 1 {
				esOpts.port = mathutil.SafeUint("1" + timex.Now().DateFormat("md")) // eg: 10425
			}

			srv := rux.New(func(r *rux.Router) {})
			srv.Use(handlers.ConsoleLogger())
			srv.Any("/{all}", func(c *rux.Context) {
				data := testutil.BuildEchoReply(c.Req)

				// c.JSON(200, data)
				c.Respond(200, data, render.NewJSONIndented())
			})

			srv.Listen("127.0.0.1", mathutil.String(esOpts.port))
			return srv.Err()
		},
	}
}
