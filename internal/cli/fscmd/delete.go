package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// DeleteCmd instance
var DeleteCmd = &gcli.Command{
	Name:    "delete",
	Desc:    "hot reload serve on files modified",
	Aliases: []string{"del", "rm"},
	Config: func(c *gcli.Command) {
		// TODO regex: from (\w+)_(\w+) to $1_new_$2
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
