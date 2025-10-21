package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
)

var configOpts = struct {
	edit bool
}{}

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
	Config: func(c *gcli.Command) {
		// builtin editors: TODO
		//  - 通用: vim, helix, nvim
		//  - Linux: nano, vi
		//  - Windows: notepad
		c.BoolOpt(&configOpts.edit, "edit", "e", false, "Edit the configuration file in the default editor")
	},
	Func: func(c *gcli.Command, args []string) error {
		// Initialize load config file
		if err := config.Mgr.Init(); err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
		cfgMgr := config.Mgr
		c.Infoln("Loading config file:", cfgMgr.Config.ConfigFile())

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
	// var configSetOpts = struct {
	// 	configPath string
	// }{}

	return &gcli.Command{
		Name:    "set",
		Desc:    "Set a configuration value",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "configuration key name", true)
			c.AddArg("value", "configuration value", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()
			value := c.Arg("value").String()

			// Initialize load config file
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			cfgMgr := config.Mgr

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

			// Save the configuration TODO
			// if err := cfgMgr.SaveConfig(""); err != nil {
			// 	return fmt.Errorf("failed to save configuration: %w", err)
			// }

			fmt.Printf("Set %s=%s\n", name, value)
			return nil
		},
	}
}

// ConfigGetCmd command for getting configuration values
func ConfigGetCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "get",
		Desc:    "Get a configuration value",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "configuration key name", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// Initialize load config file
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			cfgMgr := config.Mgr
			name := c.Arg("name").String()

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
	}
}

// ConfigExportCmd command for exporting configuration
func ConfigExportCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "export",
		Desc:    "Export configuration",
		Config: func(c *gcli.Command) {
			c.AddArg("format", "export format, allow: zip, json").WithDefault("zip")
		},
		Func: func(c *gcli.Command, args []string) error {
			format := c.Arg("format").String()
			// Validate format
			if format != "zip" && format != "json" {
				return fmt.Errorf("unsupported export format: %s (use 'zip' or 'json')", format)
			}

			// Initialize load config file
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			cfgMgr := config.Mgr

			// Create exporter
			exporter := config.NewExporter(cfgMgr)
			exportPath := "xenv_config_export." + format

			// Export the configuration
			if err := exporter.Export(exportPath, format); err != nil {
				return fmt.Errorf("failed to export configuration: %w", err)
			}

			fmt.Printf("Configuration exported to: %s\n", exportPath)
			return nil
		},
	}
}

// ConfigImportCmd command for importing configuration
func ConfigImportCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "import",
		Desc:    "Import configuration from file",
		Config: func(c *gcli.Command) {
			c.AddArg("path", "path to import configuration from", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			importPath := c.Arg("path").String()

			// Initialize load config file
			if err := config.Mgr.Init(); err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			cfgMgr := config.Mgr

			// Create importer
			importer := config.NewImporter(cfgMgr)

			// Import the configuration
			if err := importer.Import(importPath); err != nil {
				return fmt.Errorf("failed to import configuration: %w", err)
			}

			// Save the imported configuration TODO
			// if err := cfgMgr.SaveConfig("configPath"); err != nil {
			// 	return fmt.Errorf("failed to save imported configuration: %w", err)
			// }

			fmt.Printf("Configuration imported from: %s\n", importPath)
			return nil
		},
	}
}
