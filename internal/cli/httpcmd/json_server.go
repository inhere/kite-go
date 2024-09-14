package httpcmd

import "github.com/gookit/gcli/v3"

// NewJSONServerCmd instance
func NewJSONServerCmd() *gcli.Command {

	return &gcli.Command{
		Name:    "json-server",
		Desc:    "start an simple json http server",
		Aliases: []string{"json-serve", "json-srv", "jss"},
		Config: func(c *gcli.Command) {

		},
		Func: func(c *gcli.Command, args []string) error {
			return c.NewErr("TODO: not implement")
		},
	}
}
