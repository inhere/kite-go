package xenvcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/xenvcmd/subcmd"
)

// XEnvCmd the main xenv command
var XEnvCmd = &gcli.Command{
	Name:    "xenv",
	// Aliases: []string{"xenv"},
	Desc:   "Manage local development environments and tools, similar to mise and vfox",
	Subs: []*gcli.Command{
		subcmd.ToolsCmd,
		subcmd.UseCmd,
		subcmd.UnuseCmd,
		subcmd.EnvCmd,
		subcmd.PathCmd,
		subcmd.ConfigCmd,
		subcmd.ListCmd,
		subcmd.ShellCmd,
		subcmd.InitCmd,
	},
	// Configure the command
	Config: func(c *gcli.Command) {
		// Add global options for xenv command if needed
		c.BoolOpt(&subcmd.GlobalFlag, "global", "g", false, "Operate for global config")

		// Add any configuration here if needed
	},
	// Define the main command behavior (this is for the base xenv command)
	Func: func(c *gcli.Command, args []string) error {
		return c.ShowHelp()
	},
}
