package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// ListFilesCmd instance
var ListFilesCmd = &gcli.Command{
	Name:    "ls",
	Aliases: []string{"list"},
	Desc:    "list files or dirs like `ls` command",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
