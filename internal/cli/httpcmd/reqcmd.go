package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

var reqOpts = struct {
}{}

// SendRequestCmd instance
var SendRequestCmd = &gcli.Command{
	Name:    "send",
	Aliases: []string{"req", "curl"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {
		// todo: loop query, send topic, send by template
		// eg:
		// 	kite http send @jenkins trigger -v env=qa -v name=order
		// 	kite http send @feishu bot-notify
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
