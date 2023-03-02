package swagcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var InstallSwagUI = &gcli.Command{
	Name:    "downui",
	Aliases: []string{"inui", "installui"},
	Desc:    "download latest swagger-UI assets from github repository",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
