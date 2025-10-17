package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
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
		c.Infoln("Loading config file: %s", configPath)
		// Try to load existing config, ignore errors (will use defaults)
		err := cfgMgr.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Display current configuration
		fmt.Println("Current xenv configuration:")
		fmt.Printf("  BinDir: %s\n", cfgMgr.Config.BinDir)
		fmt.Printf("  InstallDir: %s\n", cfgMgr.Config.InstallDir)
		fmt.Printf("  ShellHooksDir: %s\n", cfgMgr.Config.ShellHooksDir)
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
			case "shell_hooks_dir":
				cfgMgr.Config.ShellHooksDir = value
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
			case "shell_hooks_dir":
				value = cfgMgr.Config.ShellHooksDir
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
