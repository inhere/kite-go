package sqlcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var Conv2Struct = &gcli.Command{
	Name:    "struct",
	Aliases: []string{"tostruct"},
	Desc:    "convert create table SQL to Go struct",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
