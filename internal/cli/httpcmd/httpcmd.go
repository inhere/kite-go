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
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}

// SendRequestCmd instance
var SendRequestCmd = &gcli.Command{
	Name:    "send",
	Aliases: []string{"req", "curl"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
