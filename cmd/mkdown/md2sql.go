package mkdown

import (
	"errors"

	"github.com/gookit/gcli/v2"
)

var Markdown2SQL = &gcli.Command{
	Name:    "md2sql",
	Aliases: []string{"mkdown2sql"},
	UseFor:  "convert an markdown table to create DB table SQL",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
