package doctool

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// DocumentCmd instance
var DocumentCmd = &gcli.Command{
	Name:    "doc",
	Desc:    "how to use for common tools or common commands",
	Aliases: []string{"htu"},
	Subs: []*gcli.Command{
		LinuxCmd,
		InstallCmd,
	},
}

// InstallCmd instance
var InstallCmd = &gcli.Command{
	Name:    "install",
	Aliases: []string{"ins", "add"},
	Desc:    "install new documents from git repository(eg: github)",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
