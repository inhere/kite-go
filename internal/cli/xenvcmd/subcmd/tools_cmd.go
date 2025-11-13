package subcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/xenv"
)

// ToolsCmd the xenv tools command
var ToolsCmd = &gcli.Command{
	Name:    "tools",
	Desc: "Manage local development SDK, tools (install, list, etc.)",
	Aliases: []string{"t", "tool"},
	Subs: []*gcli.Command{
		ToolsInstallCmd(),
		ToolsUninstallCmd(),
		ToolsUpdateCmd(),
		ToolsShowCmd(),
		ToolsListCmd(),
		ToolsRegisterCmd(),
		ToolsIndexCmd(),
	},
	Config: func(c *gcli.Command) {
		// Add configuration for tools command if needed
	},
}

// ToolsIndexCmd command for 将本地的工具信息索引到 local.json
func ToolsIndexCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "index",
		Help:    "index",
		Desc:    "Index local installed tools to metadata",
		Aliases: []string{"idx"},
		Func: func(c *gcli.Command, args []string) error {
			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// Index local tools
			return toolSvc.IndexLocalTools()
		},
	}

}

func ToolsRegisterCmd() *gcli.Command {
	var toolRegisterOpts = struct {
		Name    string `flag:"desc=Name of the tool to register"`
		Version string `flag:"shorts=v;desc=Version of the tool to register"`
		URL     string `flag:"name=url;desc=URL of the tool to register"`
		Bin     string `flag:"desc=Bin path of the tool to register"`
		Refresh bool   `flag:"shorts=r;desc=Refresh register tools metadata"`
	}{}

	return &gcli.Command{
		Name:    "register",
		Help:    "Register a tool",
		Desc:    "Register a tool to xenv",
		Aliases: []string{"add", "reg"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&toolRegisterOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()
			version := c.Arg("version").String()
			url := c.Arg("url").String()
			bin := c.Arg("bin").String()

			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			return toolSvc.Register(name, version, url, bin)
		},
	}

}

// ToolsInstallCmd command for installing tools
func ToolsInstallCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "install",
		Help: "install <name:version>...",
		Desc: "Install a tool with specific version",
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

			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// Install the tool
			if err := toolSvc.InstallTool(name, version); err != nil {
				return fmt.Errorf("failed to install tool %s:%s: %w", name, version, err)
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
		Help: "uninstall <name:version>",
		Desc: "Uninstall a tool with specific version",
		Aliases: []string{"un", "rm", "remove"},
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

			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// Uninstall the tool
			if err := toolSvc.Uninstall(name, version); err != nil {
				return fmt.Errorf("failed to uninstall tool %s:%s: %w", name, version, err)
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
		Help: "update <name>...",
		Desc: "Update a tool to latest or specified version",
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

			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// Update the tool (install the new version)
			if err := toolSvc.UpdateTool(name, version); err != nil {
				return fmt.Errorf("failed to update tool %s:%s: %w", name, version, err)
			}

			c.Infof("Successfully updated %s:%s\n", name, version)
			return nil
		},
	}
}

// ToolsShowCmd command for showing tool info
func ToolsShowCmd() *gcli.Command {
	return &gcli.Command{
		Name: "show",
		Help: "show <name>",
		Desc: "Show information about a specific tool",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "Name of the tool to show", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := args[0]

			// Create tool service
			toolSvc, err := xenv.ToolService()
			if err != nil {
				return err
			}

			// Get tool info
			tool := toolSvc.GetTool(name)
			if tool == nil {
				return fmt.Errorf("tool %s is not installed", name)
			}

			c.Infof("Tool: %s\n", tool.Name)
			c.Infof("  InstallDir: %s\n", tool.InstallDir)
			// c.Infof("  Installed: %t\n", tool.Installed)
			if len(tool.Alias) > 0 {
				c.Infoln(fmt.Sprintf("  Aliases: %v", tool.Alias))
			}
			return nil
		},
	}
}

// ToolsListCmd command for listing tools
func ToolsListCmd() *gcli.Command {
	var toolsListOpts = struct {
		All bool `flag:"shorts=a;desc=List all configuration tools, including uninstalled ones"`
	}{}

	return &gcli.Command{
		Name:    "list",
		Desc: "List all configuration tools",
		Aliases: []string{"ls"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&toolsListOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			return listTools()
		},
	}
}

func listTools() error {
	// Create tool service
	toolSvc, err := xenv.ToolService()
	if err != nil {
		return err
	}

	// List all tools
	return toolSvc.ListAll(false)
}
