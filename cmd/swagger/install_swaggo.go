package swagger

import (
	"errors"

	"github.com/gookit/gcli/v2"
)

var InstallSwagGo = &gcli.Command{
	Name:    "swaggo:install",
	Aliases: []string{"swaggo-install"},
	UseFor:  "install swaggo/swag from github repository",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
