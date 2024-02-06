package textcmd

import "github.com/gookit/gcli/v3"

// NewProcessCmd create a new ProcessCmd instance
func NewProcessCmd() *gcli.Command {
	var procOpts = struct {
	}{}

	return &gcli.Command{
		Name:    "process",
		Desc:    "Process input text contents",
		Aliases: []string{"proc", "handle"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&procOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			// TODO
			return nil
		},
	}
}
