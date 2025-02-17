package toolcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var XFileCmd = &gcli.Command{
	Name:    "xfile",
	Aliases: []string{"xrun"},
	Desc:    "execute kite xfile in workdir or parent dir",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
