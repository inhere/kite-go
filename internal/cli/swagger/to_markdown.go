package swagger

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var Doc2MkDown = &gcli.Command{
	Name:    "md",
	Aliases: []string{"tomd", "mkdown", "markdown"},
	Desc:    "convert swagger document file to markdown",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
