package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var lfOpt = struct {
	Dir  string `flag:"desc=the find directory, multi by comma;shorts=d"`
	Type string `flag:"desc=the find type, allow: file, dir;shorts=t"`
	Exec string `flag:"desc=execute command for each file;shorts=x"`
}{}

// ListFilesCmd instance
var ListFilesCmd = &gcli.Command{
	Name:    "ls",
	Aliases: []string{"list"},
	Desc:    "list files or dirs like `ls` command",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&lfOpt)
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
