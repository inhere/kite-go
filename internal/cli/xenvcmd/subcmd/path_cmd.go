package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv"
)

// PathCmd the xenv path command
var PathCmd = &gcli.Command{
	Name: "path",
	Desc: "Manage PATH environment variable",
	Subs: []*gcli.Command{
		PathAddCmd(),
		PathRemoveCmd(),
		PathListCmd(),
		PathSearchCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		return listEnvPaths()
	},
}

// PathAddCmd command for adding a path to PATH
func PathAddCmd() *gcli.Command {
	return &gcli.Command{
		Name: "add",
		Help: "add [-g] <path>",
		Desc: "Add a path to PATH environment variable",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.AddArg("path", "PATH environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create env service
			envSvc, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Add the path
			path := c.Arg("path").String()
			script, err1 := envSvc.AddPath(path, GlobalFlag)
			if err1 != nil {
				return fmt.Errorf("failed to add path: %w", err1)
			}

			// Save configuration if global
			if GlobalFlag {
				fmt.Printf("Added %s to PATH globally\n", path)
			} else {
				fmt.Printf("Added %s to PATH for current session\n", path)
			}

			if script != "" {
				fmt.Printf("%s\n%s\n", xenv.ScriptMark, script)
			}
			return nil
		},
	}
}

// PathRemoveCmd command for removing a path from PATH
func PathRemoveCmd() *gcli.Command {
	var pathRmOpts = struct {
		matchMode bool // TODO
	}{}

	return &gcli.Command{
		Name:    "remove",
		Help:    "remove [-g] <path>",
		Desc:    "Remove a path from PATH environment variable",
		Aliases: []string{"rm", "delete"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.BoolOpt(&pathRmOpts.matchMode, "match", "m", false, "Match mode, remove paths that match the given path")
			c.AddArg("path", "PATH environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Create env service
			envSvc, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Remove the path
			path := c.Arg("path").String()
			script, err1 := envSvc.RemovePath(path, GlobalFlag)
			if err1 != nil {
				return fmt.Errorf("failed to remove path: %w", err1)
			}

			// Save configuration if global
			if GlobalFlag {
				fmt.Printf("Removed %s from PATH globally\n", path)
			} else {
				fmt.Printf("Removed %s from PATH for current session\n", path)
			}

			if script != "" {
				fmt.Printf("%s\n%s\n", xenv.ScriptMark, script)
			}
			return nil
		},
	}
}

// PathListCmd command for listing PATH entries
func PathListCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "list",
		Desc: "List PATH entries",
		Aliases: []string{"ls"},
		Func: func(c *gcli.Command, args []string) error {
			return listEnvPaths()
		},
	}
}

func listEnvPaths() error {
	// Create env service
	envSvc, err := xenv.EnvService()
	if err != nil {
		return err
	}

	// List PATH entries
	ccolor.Infoln("Global PATH Entries:")
	for i, path := range envSvc.GlobalState().Paths {
		fmt.Printf("  %d. %s\n", i+1, path)
	}

	ccolor.Infoln("Session PATH Entries:")
	for i, path := range envSvc.SessionState().Paths {
		fmt.Printf("  %d. %s\n", i+1, path)
	}

	return nil
}

// PathSearchCmd command for searching PATH entries
func PathSearchCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "search",
		Desc: "Search for a path in PATH",
		Aliases: []string{"s"},
		Config: func(c *gcli.Command) {
			c.AddArg("value", "value for search in PATH", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			searchTerm := c.Arg("value").String()

			// Create env service
			envSvc, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Search for the path
			matches := envSvc.SearchPath(searchTerm)
			if len(matches) == 0 {
				fmt.Printf("No paths found containing: %s\n", searchTerm)
			} else {
				fmt.Printf("Paths containing '%s':\n", searchTerm)
				for i, match := range matches {
					fmt.Printf("  %d. %s\n", i+1, match)
				}
			}

			return nil
		},
	}
}
