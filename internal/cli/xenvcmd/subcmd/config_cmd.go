package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/env"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/tools"
)

// ConfigCmd the xenv config command
var ConfigCmd = &gcli.Command{
	Name:    "config",
	Desc:    "Manage xenv configuration",
	Aliases: []string{"cfg"},
	Subs: []*gcli.Command{
		ConfigSetCmd(),
		ConfigGetCmd(),
		ConfigExportCmd(),
		ConfigImportCmd(),
	},
	Func: func(c *gcli.Command, args []string) error {
		// Initialize configuration
		cfgMgr := config.NewConfigManager()
		configPath := config.GetDefaultConfigPath()
		// Try to load existing config, ignore errors (will use defaults)
		_ = cfgMgr.LoadConfig(configPath)

		// Display current configuration
		fmt.Println("Current xenv configuration:")
		fmt.Printf("  BinDir: %s\n", cfgMgr.Config.BinDir)
		fmt.Printf("  InstallDir: %s\n", cfgMgr.Config.InstallDir)
		fmt.Printf("  ShellScriptsDir: %s\n", cfgMgr.Config.ShellScriptsDir)
		fmt.Printf("  Number of managed tools: %d\n", len(cfgMgr.Config.Tools))

		return nil
	},
}

// ConfigSetCmd command for setting configuration values
func ConfigSetCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "set",
		Desc:    "Set a configuration value",
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]
			value := args[1]

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Set the configuration value based on the name
			switch name {
			case "bin_dir":
				cfgMgr.Config.BinDir = value
			case "install_dir":
				cfgMgr.Config.InstallDir = value
			case "shell_scripts_dir":
				cfgMgr.Config.ShellScriptsDir = value
			default:
				return fmt.Errorf("unknown configuration option: %s", name)
			}

			// Save the configuration
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			fmt.Printf("Set %s=%s\n", name, value)
			return nil
		},
		Config: func(c *gcli.Command) {
			c.AddArg("name", "configuration key name", true)
			c.AddArg("value", "configuration value", true)
		},
	}
}

// ConfigGetCmd command for getting configuration values
func ConfigGetCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "get",
		Desc:    "Get a configuration value",
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Get the configuration value based on the name
			var value string
			switch name {
			case "bin_dir":
				value = cfgMgr.Config.BinDir
			case "install_dir":
				value = cfgMgr.Config.InstallDir
			case "shell_scripts_dir":
				value = cfgMgr.Config.ShellScriptsDir
			default:
				return fmt.Errorf("unknown configuration option: %s", name)
			}

			fmt.Printf("%s=%s\n", name, value)
			return nil
		},
		Config: func(c *gcli.Command) {
			c.AddArg("name", "configuration key name", true)
		},
	}
}

// ConfigExportCmd command for exporting configuration
func ConfigExportCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "export",
		Desc:    "Export configuration",
		Func: func(c *gcli.Command, args []string) error {
			format := "zip" // default format
			if len(args) > 0 {
				format = args[0]
			}

			// Validate format
			if format != "zip" && format != "json" {
				return fmt.Errorf("unsupported export format: %s (use 'zip' or 'json')", format)
			}

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create exporter
			exporter := config.NewExporter(cfgMgr)

			// Determine export file path
			exportPath := "xenv_config_export." + format

			// Export the configuration
			if err := exporter.Export(exportPath, format); err != nil {
				return fmt.Errorf("failed to export configuration: %w", err)
			}

			fmt.Printf("Configuration exported to: %s\n", exportPath)
			return nil
		},
		Config: func(c *gcli.Command) {
			c.AddArg("format", "export format (zip or json)")
		},
	}
}

// ConfigImportCmd command for importing configuration
func ConfigImportCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "import",
		Desc:    "Import configuration from file",
		Func: func(c *gcli.Command, args []string) error {
			importPath := args[0]

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create importer
			importer := config.NewImporter(cfgMgr)

			// Import the configuration
			if err := importer.Import(importPath); err != nil {
				return fmt.Errorf("failed to import configuration: %w", err)
			}

			// Save the imported configuration
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save imported configuration: %w", err)
			}

			fmt.Printf("Configuration imported from: %s\n", importPath)
			return nil
		},
		Config: func(c *gcli.Command) {
			c.AddArg("path", "path to import configuration from", true)
		},
	}
}

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
