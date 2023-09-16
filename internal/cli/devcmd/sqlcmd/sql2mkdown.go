package sqlcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var Create2Mkdown = &gcli.Command{
	Name:    "to-md",
	Aliases: []string{"tomd", "to-markdown"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}

// NewCreate2JSONCmd convert create table SQL to JSON object
func NewCreate2JSONCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "c2json",
		Desc:    "parse the create table SQL to JSON object",
		Aliases: []string{"c2map", "cjson"},
		Config: func(c *gcli.Command) {

		},
		Func: func(c *gcli.Command, _ []string) error {
			return errors.New("TODO")
		},
	}
}

// Conv2StructCmd convert create table SQL to Go struct
var Conv2StructCmd = &gcli.Command{
	Name:    "struct",
	Aliases: []string{"to-struct", "tostruct", "go-struct"},
	Desc:    "convert create table SQL to Go struct",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
