package comtool

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var FileCat = &gcli.Command{
	Name:    "cat",
	Aliases: []string{"see", "bat"},
	Desc:    "hot reload serve on files modified",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
