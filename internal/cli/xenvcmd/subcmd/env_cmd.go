package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
)

// GlobalFlag option value
var GlobalFlag bool

// EnvCmd the xenv env command
var EnvCmd = &gcli.Command{
	Name: "env",
	Desc: "Manage environment variables",
	Subs: []*gcli.Command{
		EnvSetCmd(),
		EnvUnsetCmd(),
		EnvListCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		return listEnvs()
	},
}

// EnvSetCmd command for setting environment variables
func EnvSetCmd() *gcli.Command {
	return &gcli.Command{
		Name: "set",
		Help: "set [-g] <name> <value>",
		Desc: "Set an environment variable",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")

			c.AddArg("name", "environment key name", true)
			c.AddArg("value", "environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]
			value := args[1]

			// Create env manager
			envMgr, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Set the environment variable
			if err := envMgr.SetEnv(name, value, GlobalFlag); err != nil {
				return fmt.Errorf("failed to set environment variable: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				fmt.Printf("Set %s=%s globally\n", name, value)
			} else {
				fmt.Printf("Set %s=%s for current session\n", name, value)
			}

			return nil
		},
	}
}

// EnvUnsetCmd command for unsetting environment variables
func EnvUnsetCmd() *gcli.Command {
	return &gcli.Command{
		Name: "unset",
		Help: "unset [-g] <name...>",
		Desc: "Unset environment variables",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")
			c.AddArg("name", "environment key name", true)
		},
		Func: func(c *gcli.Command, args []string) error {

			// Create env manager
			envMgr, err := xenv.EnvService()
			if err != nil {
				return err
			}

			for _, name := range args {
				// Unset the environment variable
				if err := envMgr.UnsetEnv(name, GlobalFlag); err != nil {
					return fmt.Errorf("failed to unset environment variable %s: %w", name, err)
				}

				// Save configuration if global
				if GlobalFlag {
					fmt.Printf("Unset %s globally\n", name)
				} else {
					fmt.Printf("Unset %s for current session\n", name)
				}
			}

			return nil
		},
	}
}

// EnvListCmd command for listing environment variables
func EnvListCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "list",
		Desc: "List environment variables",
		Aliases: []string{"ls"},
		Func: func(c *gcli.Command, args []string) error {
			return listEnvs()
		},
	}
}

func listEnvs() error {
	// Create env manager
	envMgr, err := xenv.EnvService()
	if err != nil {
		return err
	}

	// List environment variables
	envVars := envMgr.ListEnv()
	fmt.Println("Environment Variables:")
	for name, envVar := range envVars {
		fmt.Printf("  %s=%s\n", name, envVar)
	}
	return nil
}

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
			path := c.Arg("path").String()

			// Create env manager
			envMgr, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Add the path
			if err := envMgr.AddPath(path, GlobalFlag); err != nil {
				return fmt.Errorf("failed to add path: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				fmt.Printf("Added %s to PATH globally\n", path)
			} else {
				fmt.Printf("Added %s to PATH for current session\n", path)
			}

			return nil
		},
	}
}

// PathRemoveCmd command for removing a path from PATH
func PathRemoveCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "remove",
		Help:    "remove [-g] <path>",
		Desc:    "Remove a path from PATH environment variable",
		Aliases: []string{"rm", "delete"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.AddArg("path", "PATH environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			path := c.Arg("path").String()

			// Create env manager
			envMgr, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Remove the path
			if err := envMgr.RemovePath(path, GlobalFlag); err != nil {
				return fmt.Errorf("failed to remove path: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				fmt.Printf("Removed %s from PATH globally\n", path)
			} else {
				fmt.Printf("Removed %s from PATH for current session\n", path)
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
	// Create env manager
	envMgr, err := xenv.EnvService()
	if err != nil {
		return err
	}

	// List PATH entries
	paths := envMgr.ListPaths()
	fmt.Println("PATH Entries:")
	for i, path := range paths {
		fmt.Printf("  %d. %s (%s)\n", i+1, path.Path, path.Scope)
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

			// Create env manager
			envMgr, err := xenv.EnvService()
			if err != nil {
				return err
			}

			// Search for the path
			matches := envMgr.SearchPath(searchTerm)
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
