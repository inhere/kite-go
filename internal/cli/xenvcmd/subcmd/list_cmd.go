package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/env"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/tools"
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
		cmd := ListToolsCmd()
		return cmd.Func(c, args)
	},
}

// ListToolsCmd lists tools (similar to the one in tools subcommand)
func ListToolsCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "tools",
		Desc:    "List installed tools",
		Aliases: []string{"t"},
		Func: func(c *gcli.Command, args []string) error {
			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create tool service
			toolSvc := tools.NewToolService(cfgMgr.Config)
			list := tools.NewList(toolSvc)

			// List all tools
			list.ListAll(false)

			return nil
		},
	}
}

// ListEnvCmd lists environment variables
func ListEnvCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "env",
		Desc:   "List environment variables",
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
}

// ListPathCmd lists PATH entries
func ListPathCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "path",
		Desc:   "List PATH entries",
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
}

// ListActivityCmd lists active tools and settings
func ListActivityCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "activity",
		Desc:   "List active tools and settings",
		Func: func(c *gcli.Command, args []string) error {
			// Load activity state
			activityState, err := models.LoadActivityState()
			if err != nil {
				return fmt.Errorf("failed to load activity state: %w", err)
			}

			fmt.Println("Active Tools:")
			for name, version := range activityState.ActiveTools {
				fmt.Printf("  %s:%s\n", name, version)
			}

			fmt.Println("\nActive Environment Variables:")
			for name, value := range activityState.ActiveEnv {
				fmt.Printf("  %s=%s\n", name, value)
			}

			fmt.Println("\nActive PATH Entries:")
			for i, path := range activityState.ActivePaths {
				fmt.Printf("  %d. %s\n", i+1, path)
			}

			return nil
		},
	}
}

// ListAllCmd lists everything
func ListAllCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "all",
		Desc:   "List all tools, env vars, and paths",
		Func: func(c *gcli.Command, args []string) error {
			// This would call all the other list commands
			fmt.Println("This would list all items - implementation needed")
			return nil
		},
	}
}
