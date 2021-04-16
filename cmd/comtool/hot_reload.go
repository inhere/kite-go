package comtool

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var HotReloadServe = &gcli.Command{
	Name:    "hotreload",
	Aliases: []string{"hot-reload", "hotr"},
	Desc:  "hot reload serve on files modified",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
