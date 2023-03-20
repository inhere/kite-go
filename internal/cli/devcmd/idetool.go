package devcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/inhere/kite-go/pkg/idetool"
)

// IDEAToolCmd instance
var IDEAToolCmd = &gcli.Command{
	Name:    "idea",
	Aliases: []string{"ide", "ide-tool", "jb"},
	Desc:    "find and list jetBrains tools information",
	Config: func(c *gcli.Command) {

	},
	Subs: []*gcli.Command{
		IDEAListCmd,
		IDEConfigCmd,
	},
}

// IDEAListCmd instance
var IDEAListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list installed tools on local",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		jb := idetool.NewJetBrains()

		show.AList("installed tools", jb.Installed())
		return nil
	},
}

// IDEConfigCmd instance
var IDEConfigCmd = &gcli.Command{
	Name:    "config",
	Aliases: []string{"cfg", "options"},
	Desc:    "display some tool config",
	Config: func(c *gcli.Command) {
		// search TODO
	},
	Func: func(c *gcli.Command, _ []string) error {

		return errors.New("TODO")
	},
}
