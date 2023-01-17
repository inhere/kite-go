package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// FileFinder command
var FileFinder = &gcli.Command{
	Name:    "find",
	Desc:    "hot reload serve on files modified",
	Aliases: []string{"glob"},
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
