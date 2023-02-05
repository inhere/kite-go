package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// FileCatCmd instance
var FileCatCmd = &gcli.Command{
	Name:    "cat",
	Aliases: []string{"see", "bat"},
	Desc:    "cat file contents like `ls` command",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
