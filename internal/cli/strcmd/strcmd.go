package strcmd

import "github.com/gookit/gcli/v3"

var StringCmd = &gcli.Command{
	Name:    "text",
	Desc:    "useful commands for handle string text",
	Aliases: []string{"str", "string"},
	Subs:    []*gcli.Command{},
}
