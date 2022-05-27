package mkdown

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var Markdown2SQL = &gcli.Command{
	Name:    "sql",
	Aliases: []string{"tosql"},
	Desc:    "convert an markdown table to create DB table SQL",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
