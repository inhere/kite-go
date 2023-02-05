package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/toolcmd/doctool"
	"github.com/inhere/kite/internal/cli/toolcmd/mdcmd"
	"github.com/inhere/kite/internal/cli/toolcmd/swagger"
)

// ToolsCmd command
var ToolsCmd = &gcli.Command{
	Name:    "tool",
	Aliases: []string{"tools"},
	Desc:    "provide some useful tools commands",
	Subs: []*gcli.Command{
		swagger.SwaggerCmd,
		BatchRunCmd,
		EnvInfoCmd,
		AutoJumpCmd,
		RunAnyCmd,
		doctool.DocumentCmd,
		mdcmd.MkDownCmd,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
