package subcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/inhere/kite-go/internal/util"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

var GlobalFlag bool

var shellCmdOpts = struct {
	Type gflag.EnumString
}{
	Type: cflag.NewEnumString("bash", "zsh", "pwsh"),
}

// ShellCmd the xenv shell command
var ShellCmd = &gcli.Command{
	Name: "shell",
	Desc: "Generate shell integration script",
	Func: func(c *gcli.Command, args []string) error {
		shellType := shellCmdOpts.Type.String()
		if shellType == "" {
			shellType = "bash" // default
		}

		var hookScript string
		switch shellType {
		case "bash":
			hookScript = shell.GenerateBashHook()
		case "zsh":
			hookScript = shell.GenerateZshHook()
		case "pwsh":
			hookScript = shell.GeneratePwshHook()
		default:
			return fmt.Errorf("unsupported shell type: %s (use bash, zsh, or pwsh)", shellType)
		}

		fmt.Print(hookScript)
		return nil
	},
	Config: func(c *gcli.Command) {
		c.VarOpt(&shellCmdOpts.Type, "type", "t", "Shell type (bash, zsh, pwsh)")
	},
}

// InitCmd the xenv init command
var InitCmd = &gcli.Command{
	Name: "init",
	Desc: "Initialize xenv configuration and environment",
	Func: func(c *gcli.Command, args []string) error {
		// Initialize configuration
		cfgMgr := config.NewConfigManager()
		configPath := config.GetDefaultConfigPath()

		// Ensure config directory exists
		configDir := filepath.Dir(configPath)
		if err := util.EnsureDir(configDir); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Try to load existing config, if it exists
		if _, err := os.Stat(configPath); err == nil {
			if err := cfgMgr.LoadConfig(configPath); err != nil {
				fmt.Printf("Warning: failed to load existing config: %v\n", err)
			}
		} else {
			// If no existing config, save the default config
			if err := cfgMgr.SaveConfig(configPath); err != nil {
				return fmt.Errorf("failed to save default config: %w", err)
			}
			fmt.Printf("Created default configuration at: %s\n", configPath)
		}

		// Ensure required directories exist
		if err := util.EnsureDir(util.ExpandHome(cfgMgr.Config.BinDir)); err != nil {
			return fmt.Errorf("failed to create bin directory: %w", err)
		}

		if err := util.EnsureDir(util.ExpandHome(cfgMgr.Config.InstallDir)); err != nil {
			return fmt.Errorf("failed to create install directory: %w", err)
		}

		if err := util.EnsureDir(util.ExpandHome(cfgMgr.Config.ShellHooksDir)); err != nil {
			return fmt.Errorf("failed to create shell scripts directory: %w", err)
		}

		fmt.Println("xenv initialization completed successfully!")
		fmt.Printf("Configuration file: %s\n", configPath)
		fmt.Printf("Bin directory: %s\n", cfgMgr.Config.BinDir)
		fmt.Printf("Install directory: %s\n", cfgMgr.Config.InstallDir)
		fmt.Printf("Shell scripts directory: %s\n", cfgMgr.Config.ShellHooksDir)

		return nil
	},
}
