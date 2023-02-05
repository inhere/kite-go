package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// RenameCmd instance
var RenameCmd = &gcli.Command{
	Name: "rename",
	Desc: "hot reload serve on files modified",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
