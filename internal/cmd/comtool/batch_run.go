package comtool

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// BatchRun command
var BatchRun = &gcli.Command{
	Name:    "brun",
	Aliases: []string{"batch-run"},
	Desc:    "batch run more commands at once",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
