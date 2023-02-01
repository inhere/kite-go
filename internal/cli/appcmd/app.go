package appcmd

import (
	"github.com/gookit/gcli/v3"
)

// ManageCmd manage kite self
var ManageCmd = &gcli.Command{
	Name:    "app",
	Aliases: []string{"self"},
	Desc:    "provide commands for manage kite self",
	Subs: []*gcli.Command{
		AppCheckCmd,
		KiteInitCmd,
		KiteInfoCmd,
		UpdateSelfCmd,
		KiteConfCmd,
		ReadmeCmd,
		KitePathCmd,
	},
}

// UpdateSelfCmd command
var UpdateSelfCmd = &gcli.Command{
	Name:    "update",
	Aliases: []string{"update-self", "up-self", "up"},
	Desc:    "update {$binName} to latest from github repository",
}
