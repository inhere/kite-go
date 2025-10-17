package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/env"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

// EnvCmd the xenv env command
var EnvCmd = &gcli.Command{
	Name:    "env",
	Desc:   "Manage environment variables",
	Subs: []*gcli.Command{
		EnvSetCmd(),
		EnvUnsetCmd(),
		EnvListCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		// Initialize configuration
		cfgMgr := config.NewConfigManager()
		configPath := config.GetDefaultConfigPath()
		// Try to load existing config, ignore errors (will use defaults)
		_ = cfgMgr.LoadConfig(configPath)

		// Load activity state
		activityState, err := models.LoadActivityState()
		if err != nil {
			return fmt.Errorf("failed to load activity state: %w", err)
		}

		// Create env manager
		envMgr := env.NewManager(cfgMgr.Config, activityState)

		// List environment variables
		envVars := envMgr.ListEnv()
		fmt.Println("Environment Variables:")
		for name, envVar := range envVars {
			fmt.Printf("  %s=%s (%s)\n", name, envVar.Value, envVar.Scope)
		}

		// Also show session variables from activity state
		if len(activityState.ActiveEnv) > 0 {
			fmt.Println("\nSession Variables:")
			for name, value := range activityState.ActiveEnv {
				fmt.Printf("  %s=%s\n", name, value)
			}
		}

		return nil
	},
}

// EnvSetCmd command for setting environment variables
func EnvSetCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "set",
		Help:     "set [-g] <name> <value>",
		Desc:   "Set an environment variable",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")

			c.AddArg("name", "environment key name", true)
			c.AddArg("value", "environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]
			value := args[1]

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create env manager
			envMgr := env.NewManager(cfgMgr.Config, activityState)

			// Set the environment variable
			if err := envMgr.SetEnv(name, value, GlobalFlag); err != nil {
				return fmt.Errorf("failed to set environment variable: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				if err := cfgMgr.SaveConfig(configPath); err != nil {
					return fmt.Errorf("failed to save configuration: %w", err)
				}
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
		Name:    "unset",
		Help:     "unset [-g] <name...>",
		Desc:   "Unset environment variables",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Operate for global config")
			c.AddArg("name", "environment key name", true)
		},
		Func: func(c *gcli.Command, args []string) error {

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create env manager
			envMgr := env.NewManager(cfgMgr.Config, activityState)

			for _, name := range args {
				// Unset the environment variable
				if err := envMgr.UnsetEnv(name, GlobalFlag); err != nil {
					return fmt.Errorf("failed to unset environment variable %s: %w", name, err)
				}

				// Save configuration if global
				if GlobalFlag {
					if err := cfgMgr.SaveConfig(configPath); err != nil {
						return fmt.Errorf("failed to save configuration: %w", err)
					}
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
		Desc:   "List environment variables",
		Aliases: []string{"ls"},
		Func: func(c *gcli.Command, args []string) error {
			// This is the same as the main command's Run function
			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create env manager
			envMgr := env.NewManager(cfgMgr.Config, activityState)

			// List environment variables
			envVars := envMgr.ListEnv()
			fmt.Println("Environment Variables:")
			for name, envVar := range envVars {
				fmt.Printf("  %s=%s (%s)\n", name, envVar.Value, envVar.Scope)
			}

			// Also show session variables from activity state
			if len(activityState.ActiveEnv) > 0 {
				fmt.Println("\nSession Variables:")
				for name, value := range activityState.ActiveEnv {
					fmt.Printf("  %s=%s\n", name, value)
				}
			}

			return nil
		},
	}
}

// PathCmd the xenv path command
var PathCmd = &gcli.Command{
	Name:    "path",
	Desc:   "Manage PATH environment variable",
	Subs: []*gcli.Command{
		PathAddCmd(),
		PathRemoveCmd(),
		PathListCmd(),
		PathSearchCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		// Initialize configuration
		cfgMgr := config.NewConfigManager()
		configPath := config.GetDefaultConfigPath()
		// Try to load existing config, ignore errors (will use defaults)
		_ = cfgMgr.LoadConfig(configPath)

		// Load activity state
		activityState, err := models.LoadActivityState()
		if err != nil {
			return fmt.Errorf("failed to load activity state: %w", err)
		}

		// Create path manager
		pathMgr := env.NewPathManager(cfgMgr.Config, activityState)

		// List PATH entries
		paths := pathMgr.ListPaths()
		fmt.Println("PATH Entries:")
		for i, path := range paths {
			fmt.Printf("  %d. %s (%s)\n", i+1, path.Path, path.Scope)
		}

		// Also show session paths from activity state
		if len(activityState.ActivePaths) > 0 {
			fmt.Println("\nSession PATH Entries:")
			for i, path := range activityState.ActivePaths {
				fmt.Printf("  %d. %s\n", i+1, path)
			}
		}

		return nil
	},
}

// PathAddCmd command for adding a path to PATH
func PathAddCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "add",
		Help:     "add [-g] <path>",
		Desc:   "Add a path to PATH environment variable",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.AddArg("path", "PATH environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			path := c.Arg("path").String()

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create path manager
			pathMgr := env.NewPathManager(cfgMgr.Config, activityState)

			// Add the path
			if err := pathMgr.AddPath(path, GlobalFlag); err != nil {
				return fmt.Errorf("failed to add path: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				if err := cfgMgr.SaveConfig(configPath); err != nil {
					return fmt.Errorf("failed to save configuration: %w", err)
				}
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
		Name:    "rm",
		Help:     "rm [-g] <path>",
		Desc:   "Remove a path from PATH environment variable",
		Aliases: []string{"remove"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&GlobalFlag, "global", "g", false, "Global operation, not the current session")
			c.AddArg("path", "PATH environment value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			path := c.Arg("path").String()

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create path manager
			pathMgr := env.NewPathManager(cfgMgr.Config, activityState)

			// Remove the path
			if err := pathMgr.RemovePath(path, GlobalFlag); err != nil {
				return fmt.Errorf("failed to remove path: %w", err)
			}

			// Save configuration if global
			if GlobalFlag {
				if err := cfgMgr.SaveConfig(configPath); err != nil {
					return fmt.Errorf("failed to save configuration: %w", err)
			}
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
		Desc:   "List PATH entries",
		Aliases: []string{"ls"},
		Func: func(c *gcli.Command, args []string) error {
			// This is the same as the main command's Run function
			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create path manager
			pathMgr := env.NewPathManager(cfgMgr.Config, activityState)

			// List PATH entries
			paths := pathMgr.ListPaths()
			fmt.Println("PATH Entries:")
			for i, path := range paths {
				fmt.Printf("  %d. %s (%s)\n", i+1, path.Path, path.Scope)
			}

			// Also show session paths from activity state
			if len(activityState.ActivePaths) > 0 {
				fmt.Println("\nSession PATH Entries:")
				for i, path := range activityState.ActivePaths {
					fmt.Printf("  %d. %s\n", i+1, path)
				}
			}

			return nil
		},
	}
}

// PathSearchCmd command for searching PATH entries
func PathSearchCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "search",
		Desc:   "Search for a path in PATH",
		Aliases: []string{"s"},
		Config: func(c *gcli.Command) {
			c.AddArg("value", "value for search in PATH", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			searchTerm := c.Arg("value").String()

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			// Create path manager
			pathMgr := env.NewPathManager(cfgMgr.Config, activityState)

			// Search for the path
			matches := pathMgr.SearchPath(searchTerm)
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
