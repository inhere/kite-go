package gocmd

import "github.com/gookit/gcli/v3"

// GoToolsCmd instance
var GoToolsCmd = &gcli.Command{
	Name: "go",
	Desc: "some go tools command",
	Subs: []*gcli.Command{
		AwesomeGoCmd,
		GoInfoCmd,
	},
}

// GoInfoCmd refer https://github.com/lucor/goinfo
var GoInfoCmd = &gcli.Command{
	Name: "info",
	Desc: "system info for go",
}
