package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/rux"
	"github.com/gookit/rux/render"
)

// NewEchoServerCmd instance
func NewEchoServerCmd() *gcli.Command {
	var esOpts = struct {
		port uint
	}{}

	return &gcli.Command{
		Name:    "echo",
		Desc:    "start an simple echo http server",
		Aliases: []string{"echo-serve"},
		Config: func(c *gcli.Command) {
			c.UintOpt(&esOpts.port, "port", "P", 0, "custom the echo server port, default will use random `port`")
		},
		Func: func(c *gcli.Command, args []string) error {
			if esOpts.port < 1 {
				esOpts.port = uint(mathutil.RandInt(6000, 9999))
			}

			srv := rux.New(func(r *rux.Router) {

			})

			srv.Any("/{all}", func(c *rux.Context) {
				bs, err := c.RawBodyData()
				if err != nil {
					c.AbortThen().AddError(err)
					return
				}

				data := rux.M{
					"headers": c.Req.Header,
					"uri":     c.Req.RequestURI,
					"query":   c.QueryValues(),
					"body":    string(bs),
				}

				// c.JSON(200, data)
				c.Respond(200, data, render.NewJSONIndented())
			})

			srv.Listen(mathutil.String(esOpts.port))
			return srv.Err()
		},
	}
}