package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// ListCmd the xenv list command
var ListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "List installed tools, environment variables, or PATH entries",
	Aliases: []string{"ls"},
	Subs: []*gcli.Command{
		ListToolsCmd(),
		ListEnvCmd(),
		ListPathCmd(),
		ListActivityCmd(),
		ListAllCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		// Default to listing tools if no subcommand is specified
		return listTools()
	},
}

// ListToolsCmd lists tools (similar to the one in tools subcommand)
func ListToolsCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "tools",
		Desc: "List local installed tools",
		Aliases: []string{"t"},
		Func: func(c *gcli.Command, args []string) error {
			return listTools()
		},
	}
}

// ListEnvCmd lists environment variables
func ListEnvCmd() *gcli.Command {
	return &gcli.Command{
		Name: "env",
		Desc: "List environment variables",
		Func: func(c *gcli.Command, args []string) error {
			return listEnvs()
		},
	}
}

// ListPathCmd lists PATH entries
func ListPathCmd() *gcli.Command {
	return &gcli.Command{
		Name: "path",
		Desc: "List PATH entries",
		Func: func(c *gcli.Command, args []string) error {
			return listEnvPaths()
		},
	}
}

// ListActivityCmd lists active tools and settings
func ListActivityCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "activity",
		Desc:    "List active tools and settings",
		Aliases: []string{"act", "active", "use"},
		Func: func(c *gcli.Command, args []string) error {
			// Load activity state
			if err := xenv.State().Init(); err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			show.ATitle("Global State")
			global := xenv.State().Global()
			if global.IsEmpty() {
				ccolor.Infoln("No global state found")
			} else {
				listActivity(global)
			}

			dirStates := xenv.State().DirStates()
			if len(dirStates) > 0 {
				fmt.Println()
				show.ATitle("Directory States")
				for _, dirState := range dirStates {
					ccolor.Infoln(" - form: %s", dirState.File)
					listActivity(dirState)
				}
			}

			if util.InHookShell() {
				fmt.Println()
				show.ATitle("Session State")
				listActivity(xenv.State().Session())
			}
			return nil
		},
	}
}

func listActivity(state *models.ActivityState) {
	ccolor.Cyanln("Active SDK Tools:")
	for name, version := range state.SDKs {
		ccolor.Printf("  <green>%s</> => %s\n", name, version)
	}

	ccolor.Cyanln("\nActive Environment Variables:")
	for name, value := range state.Envs {
		ccolor.Printf("  <green>%s</>=%s\n", name, value)
	}

	ccolor.Cyanln("\nActive PATH Entries:")
	for i, path := range state.Paths {
		ccolor.Printf("  <green>%d</>. %s\n", i+1, path)
	}
}

// ListAllCmd lists everything
func ListAllCmd() *gcli.Command {
	return &gcli.Command{
		Name: "all",
		Desc: "List all tools, env vars, and paths",
		Func: func(c *gcli.Command, args []string) error {
			// This would call all the other list commands
			fmt.Println("This would list all items - implementation needed")
			return nil
		},
	}
}
