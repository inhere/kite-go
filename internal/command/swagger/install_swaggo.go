package swagger

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var InstallSwagGo = &gcli.Command{
	Name:    "swaggo",
	Aliases: []string{"swaggo-ins"},
	Desc:    "install swaggo/swag from github repository",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
