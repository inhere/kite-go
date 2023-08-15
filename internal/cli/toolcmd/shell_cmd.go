package toolcmd

import (
	"github.com/gookit/gcli/v3"
)

// ShellCmd os shell tool command
var ShellCmd = &gcli.Command{
	Name:    "shell",
	Aliases: []string{"sh"},
	Desc:    "Listen, record, query user executed shell command",
	Subs: []*gcli.Command{
		ShellListenCmd,
	},
	Config: func(c *gcli.Command) {
		// c.AddArg("name", "The name of user")
	},
}

// ShellListenCmd listen user executed command
var ShellListenCmd = &gcli.Command{
	Name:    "listen",
	Aliases: []string{"watch"},
	Desc:    "Listen user executed shell command",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, args []string) error {
		return nil
	},
}
