package swagger

import (
	"errors"

	"github.com/gookit/gcli/v2"
)

var Doc2MkDown = &gcli.Command{
	Name:    "swag2md",
	Aliases: []string{"swagtomd", "swag:tomd"},
	UseFor:  "convert swagger document file to markdown",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
