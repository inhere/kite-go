package gencmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// NewJsonToCodeCmd instance
func NewJsonToCodeCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "json2code",
		Aliases: []string{"j2c"},
		Desc:    "generate java/php/go code for json(5) codes",
		Func: func(c *gcli.Command, _ []string) error {
			return errorx.New("TODO")
		},
	}
}
