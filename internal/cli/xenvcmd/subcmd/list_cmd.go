package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show/title"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
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
		Desc:    "List local installed SDK tools",
		Aliases: []string{"t", "tool", "sdk", "sdks"},
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
	var listActOpts = struct {
		Group bool `flag:"shorts=t;desc=List activity states and group by global, dir, session"`
	}{}

	return &gcli.Command{
		Name:    "activity",
		Desc: "List active SDKs, envs, paths and tools",
		Aliases: []string{"act", "active", "use"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&listActOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Load activity state
			if err := xenv.InitState(); err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			tl := title.New("", func(t *title.Title) {
				t.Width = 40
				t.PaddingLR = false
				t.ShowBorder = true
			})
			if !listActOpts.Group {
				tl.ShowNew("Activity States")
				listActivity(xenv.State().Merged())
				return nil
			}

			tl.ShowNew("Global State")
			global := xenv.State().Global()
			if global.IsEmpty() {
				ccolor.Infoln("No global state found")
			} else {
				listActivity(global)
			}

			dirStates := xenv.State().DirStates()
			if len(dirStates) > 0 {
				fmt.Println()
				tl.ShowNew("Directory States")
				for _, dirState := range dirStates {
					fmt.Println(" - form:", dirState.File)
					listActivity(dirState)
				}
			}

			if xenvcom.InHookShell() {
				fmt.Println()

				sess := xenv.State().Session()
				tl.ShowNew("Session State")
				fmt.Println(" - session ID:", sess.SessionID())
				if global.IsEmpty() {
					ccolor.Infoln("No session state found")
				} else {
					listActivity(sess)
				}
			}
			return nil
		},
	}
}

func listActivity(state *models.ActivityState) {
	ccolor.Cyanln("Active Develop SDKs:")
	for name, version := range state.SDKs {
		ccolor.Printf("  <green>%s</> => %s\n", name, version)
	}

	ccolor.Cyanln("\nActive Env Variables:")
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
