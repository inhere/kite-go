package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// HttpCmd command
var HttpCmd = &gcli.Command{
	Name: "http",
	// Aliases: []string{"h"},
	Desc: "provide some useful tools commands",
	Subs: []*gcli.Command{
		HttpServeCmd,
		SendRequestCmd,
		NewEchoServerCmd(),
		NewFileServerCmd(),
	},
	Config: func(c *gcli.Command) {

	},
}

// SendRequestCmd instance
var SendRequestCmd = &gcli.Command{
	Name:    "send",
	Aliases: []string{"req", "curl"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {
		// todo: loop query, send topic, send by template
		// eg:
		// 	kite http send @jenkins trigger
		// 	kite http send @feishu bot-notify
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
