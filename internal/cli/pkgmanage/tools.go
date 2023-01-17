package pkgmanage

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// ManageCmd command
var ManageCmd = &gcli.Command{
	Name:    "pkgm",
	Aliases: []string{"pkgx"},
	Desc:    "local tools management. eg: install, update and remove",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
