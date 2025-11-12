package xenvcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli/xenvcmd/subcmd"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

// XEnvCmd the main xenv command
var XEnvCmd = &gcli.Command{
	Name:    "xenv",
	// Aliases: []string{"xenv"},
	Desc:   "Manage local development environments and tools, similar to mise and vfox",
	Help: `
Quick commands:
  <info>set</>    Quick exec the 'env set' subcommand
  <info>unset</>  Quick exec the 'env unset' subcommand
`,
	Subs: []*gcli.Command{
		subcmd.ToolsCmd,
		subcmd.NewUseCmd(),
		subcmd.NewUnuseCmd(),
		subcmd.EnvCmd,
		subcmd.PathCmd,
		subcmd.ConfigCmd,
		subcmd.ListCmd,
		subcmd.InitCmd,
		subcmd.NewShellCmd(),
		subcmd.ShellHookInitCmd(),
		subcmd.ShellDirenvCmd(),
	},
	// Configure the command
	Config: func(c *gcli.Command) {
		// Add global options for xenv command if needed
		c.BoolOpt(&subcmd.GlobalFlag, "global", "g", false, "Operate for global config")
		c.BoolOpt(&xenvcom.DebugMode, "debug", "d", false, "Enable debug mode. can be XENV_DEBUG_MODE=true")

		// Add any configuration here if needed
		c.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) (stop bool) {
			name := ctx.Str("name")
			// 重定向执行 env set/unset 命令
			if name == "set" || name == "unset" {
				newArgs := []string{"env", name}
				newArgs = append(newArgs, ctx.Strings("args")...)
				err := app.Cli.RunCmd("xenv", newArgs)
				if err != nil {
					fmt.Println(err)
				}
				return true
			}
			return false
		})
	},
}
