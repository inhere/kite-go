package jsoncmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// JSONToolCmd instance
var JSONToolCmd = &gcli.Command{
	Name: "json",
	Desc: "json tool commands",
	Subs: []*gcli.Command{
		JSONViewCmd,
	},
}

var JSONViewCmd = &gcli.Command{
	Name:    "view",
	Aliases: []string{"cat", "fmt"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
