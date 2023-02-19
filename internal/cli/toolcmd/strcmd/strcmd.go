package strcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

var StringCmd = &gcli.Command{
	Name:    "text",
	Desc:    "useful commands for handle string text",
	Aliases: []string{"str", "string"},
	Subs: []*gcli.Command{
		StrCountCmd,
	},
}

// StrCountCmd instance
var StrCountCmd = &gcli.Command{
	Name:    "length",
	Aliases: []string{"len", "count"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
