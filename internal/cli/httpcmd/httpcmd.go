package httpcmd

import (
	"github.com/gookit/gcli/v3"
)

// HttpCmd command
var HttpCmd = &gcli.Command{
	Name: "http",
	// Aliases: []string{"h"},
	Desc: "provide some useful tools commands",
	Subs: []*gcli.Command{
		HttpServeCmd,
		SendRequestCmd,
		SendTemplateCmd,
		TemplateInfoCmd,
		DecodeQueryCmd,
		NewEchoServerCmd(),
		NewFileServerCmd(),
		NewHookServerCmd(),
		NewOAPIServeCmd(),
		// TODO convert to CURL command. refer: https://github.com/moul/http2curl/blob/master/http2curl.go
	},
	Config: func(c *gcli.Command) {

	},
}
