package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/devcmd/jsoncmd"
	"github.com/inhere/kite/internal/cli/pkgcmd"
	"github.com/inhere/kite/internal/cli/syscmd"
	"github.com/inhere/kite/internal/cli/toolcmd/doctool"
	"github.com/inhere/kite/internal/cli/toolcmd/mdcmd"
	"github.com/inhere/kite/internal/cli/toolcmd/strcmd"
	"github.com/inhere/kite/internal/cli/toolcmd/swagger"
)

// ToolsCmd command
var ToolsCmd = &gcli.Command{
	Name:    "tool",
	Aliases: []string{"tools"},
	Desc:    "provide some useful tools commands",
	Subs: []*gcli.Command{
		swagger.SwaggerCmd,
		strcmd.StringCmd,
		syscmd.NewBatchRunCmd(),
		EnvInfoCmd,
		AutoJumpCmd,
		RunAnyCmd,
		pkgcmd.PkgManageCmd,
		doctool.DocumentCmd,
		mdcmd.MkDownCmd,
		syscmd.QuickOpenCmd,
		jsoncmd.JSONToolCmd,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
