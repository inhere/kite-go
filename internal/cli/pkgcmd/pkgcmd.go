package pkgcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// PkgManageCmd command
var PkgManageCmd = &gcli.Command{
	Name:    "pkg",
	Aliases: []string{"pkgm", "pkgx"},
	Desc:    "local package tools management",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
