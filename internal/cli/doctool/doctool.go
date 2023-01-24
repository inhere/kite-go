package doctool

import "github.com/gookit/gcli/v3"

// DocumentCmd instance
var DocumentCmd = &gcli.Command{
	Name:    "doc",
	Desc:    "how to use for common tools or common commands",
	Aliases: []string{"htu"},
	Subs: []*gcli.Command{
		LinuxCmd,
	},
}
