package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/tools"
)



// ToolsCmd the xenv tools command
var ToolsCmd = &gcli.Command{
	Name:    "tools",
	Aliases: []string{"t", "tool"},
	Desc:   "Manage local development tools (install, list, etc.)",
	Subs: []*gcli.Command{
		ToolsInstallCmd(),
		ToolsUninstallCmd(),
		ToolsUpdateCmd(),
		ToolsShowCmd(),
		ToolsListCmd(),
	},
	Config: func(c *gcli.Command) {
		// Add configuration for tools command if needed
	},
	Func: func(c *gcli.Command, args []string) error {
		return c.ShowHelp()
	},
}

// ToolsInstallCmd command for installing tools
func ToolsInstallCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "install",
		Help:     "install <name:version>...",
		Desc:   "Install a tool with specific version",
		Aliases: []string{"i", "in"},
		Config: func(c *gcli.Command) {
			c.AddArg("tools", "Name of the tool to install, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			inTools := c.Arg("tools").Strings()
			// Parse name:version
			name, version, err := parseNameVersion(inTools[0])
			if err != nil {
				return err
			}

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create tool service
			toolSvc := tools.NewToolService(cfgMgr.Config)

			// Install the tool
			if err := toolSvc.InstallTool(name, version); err != nil {
				return fmt.Errorf("failed to install tool %s:%s: %w", name, version, err)
			}

			// Save configuration
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			c.Infof("Successfully installed %s:%s\n", name, version)
			return nil
		},
	}
}

// ToolsUninstallCmd command for uninstalling tools
func ToolsUninstallCmd() *gcli.Command {
	// Parse flag to determine if we should keep config
	var keepConfig bool

	return &gcli.Command{
		Name:    "uninstall",
		Help:     "uninstall <name:version>",
		Desc:   "Uninstall a tool with specific version",
		Aliases: []string{"un"},
		Config: func(c *gcli.Command) {
			// Add option to keep configuration files
			c.BoolOpt(&keepConfig, "keep-config", "kc", false, "Keep configuration files after uninstall")
		},
		Func: func(c *gcli.Command, args []string) error {
			// Parse name:version
			name, version, err := parseNameVersion(args[0])
			if err != nil {
				return err
			}

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create tool service
			toolSvc := tools.NewToolService(cfgMgr.Config)
			uninstaller := tools.NewUninstaller(toolSvc, cfgMgr.Config)

			// Uninstall the tool
			if err := uninstaller.Uninstall(name, version, keepConfig); err != nil {
				return fmt.Errorf("failed to uninstall tool %s:%s: %w", name, version, err)
			}

			// Save configuration
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			c.Infof("Successfully uninstalled %s:%s\n", name, version)
			return nil
		},
	}
}

// ToolsUpdateCmd command for updating tools
func ToolsUpdateCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "update",
		Help:     "update <name>...",
		Desc:   "Update a tool to latest or specified version",
		Aliases: []string{"up"},
		Config: func(c *gcli.Command) {
			c.AddArg("tools", "Name of the tool to update, allow multi.", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			upTools := c.Arg("tools").Strings()
			// Parse name:version
			name, version, err := parseNameVersion(upTools[0])
			if err != nil {
				return err
			}

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create tool service
			toolSvc := tools.NewToolService(cfgMgr.Config)

			// Update the tool (install the new version)
			if err := toolSvc.UpdateTool(name, version); err != nil {
				return fmt.Errorf("failed to update tool %s:%s: %w", name, version, err)
			}

			// Save configuration
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			c.Infof("Successfully updated %s:%s\n", name, version)
			return nil
		},
	}
}

// ToolsShowCmd command for showing tool info
func ToolsShowCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "show",
		Help:     "show <name>",
		Desc:   "Show information about a specific tool",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "Name of the tool to show", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]

			// Initialize configuration
			cfgMgr := config.NewConfigManager()
			configPath := config.GetDefaultConfigPath()
			// Try to load existing config, ignore errors (will use defaults)
			_ = cfgMgr.LoadConfig(configPath)

			// Create tool service
			toolSvc := tools.NewToolService(cfgMgr.Config)

			// Get tool info
			tool := toolSvc.GetTool(name)
			if tool == nil {
				return fmt.Errorf("tool %s is not installed", name)
			}

			c.Infof("Tool: %s\n", tool.ID)
			c.Infof("  InstallDir: %s\n", tool.InstallDir)
			c.Infof("  Installed: %t\n", tool.Installed)
			if len(tool.Alias) > 0 {
				c.Infoln(fmt.Sprintf("  Aliases: %v", tool.Alias))
			}
			if len(tool.BinPaths) > 0 {
				c.Infof("  BinPaths: %v\n", tool.BinPaths)
			}

			return nil
		},
	}
}

// ToolsListCmd command for listing tools
func ToolsListCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "list",
		Desc:   "List all installed tools",
		Aliases: []string{"ls"},
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
