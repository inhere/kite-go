package doctool

import "github.com/gookit/gcli/v3"

var DocumentCmd = &gcli.Command{
	Name: "doc",
	Desc: "how to use for common tools or common commands",
	Subs: []*gcli.Command{
		LinuxCmd,
	},
}
