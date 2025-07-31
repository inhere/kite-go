package appcmd

import (
	"github.com/gookit/gcli/v3"
)

// ManageCmd manage kite self
var ManageCmd = &gcli.Command{
	Name:    "app",
	Aliases: []string{"self", "mgr"},
	Desc:    "provide commands for manage kite self",
	Subs: []*gcli.Command{
		AppCheckCmd,
		KiteInitCmd,
		KiteInfoCmd,
		KiteObjectCmd,
		KiteConfCmd,
		KiteReadmeCmd,
		KitePathCmd,
		NewPathMapCmd(),
		NewAppExtCmd(),
		KiteAliasCmd,
		BackendServeCmd,
		CommandMapCmd,
		UpdateSelfCmd,
		LogWriteCmd,
	},
}

// UpdateSelfCmd command
var UpdateSelfCmd = &gcli.Command{
	Name:    "update",
	Aliases: []string{"update-self", "up-self", "up"},
	Desc:    "update {$binName} to latest from github repository",
}
