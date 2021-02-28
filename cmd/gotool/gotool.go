package gotool

import "github.com/gookit/gcli/v3"

var GoToolsCmd = &gcli.Command{
	Name: "go",
	Desc: "some go tools command",
	Subs: []*gcli.Command{AwesomeGo},
}
