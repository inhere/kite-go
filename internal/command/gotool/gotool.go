package gotool

import "github.com/gookit/gcli/v3"

var GoToolsCmd = &gcli.Command{
	Name: "go",
	Desc: "some go tools command",
	Subs: []*gcli.Command{
		AwesomeGo,
		GoInfo,
	},
}

// refer https://github.com/lucor/goinfo
var GoInfo = &gcli.Command{
	Name: "info",
	Desc: "system info for go",
}
