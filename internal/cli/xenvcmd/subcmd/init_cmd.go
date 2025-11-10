package subcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv"
	"github.com/inhere/kite-go/pkg/xenv/config"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// InitCmd the xenv init command
var InitCmd = &gcli.Command{
	Name: "init",
	Desc: "Initialize xenv configuration and environment",
	Func: func(c *gcli.Command, args []string) error {
		// Initialize configuration
		// Initialize load config file
		if err := config.Mgr.Init(); err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		cfgMgr := config.Mgr
		cfg := config.Mgr.Config
		// c.Infoln("Loading config file:", cfg.ConfigFile())
		configPath := cfg.ConfigFile()

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
		if err := util.EnsureDir(cfgMgr.Config.BinDir); err != nil {
			return fmt.Errorf("failed to create bin directory: %w", err)
		}

		if err := util.EnsureDir(cfgMgr.Config.InstallDir); err != nil {
			return fmt.Errorf("failed to create install directory: %w", err)
		}

		if err := util.EnsureDir(cfgMgr.Config.ShellHooksDir); err != nil {
			return fmt.Errorf("failed to create shell scripts directory: %w", err)
		}

		fmt.Println("Xenv initialization completed successfully!")
		fmt.Printf("Configuration file: %s\n", configPath)
		fmt.Printf("Bin directory: %s\n", cfgMgr.Config.BinDir)
		fmt.Printf("Install directory: %s\n", cfgMgr.Config.InstallDir)
		fmt.Printf("Shell scripts directory: %s\n", cfgMgr.Config.ShellHooksDir)

		return nil
	},
}

// HookInitCmd the xenv hook init command
//   - 配置了 xenv shell 命令到 user 配置文件后，会自动执行该命令
//   - 调用当前命令，可以返回脚本内容自动执行
var HookInitCmd = &gcli.Command{
	Hidden: true, // This is an internal command
	Name:   "hook-init",
	Desc:   "Initialize the xenv hook script",
	Func: func(c *gcli.Command, args []string) error {
		return nil // TODO
	},
}

// InitDirenvCmd the xenv init direnv command
//  - 仅在配置了 xenv shell 命令时，cd 到新目录会自动调用当前命令
//  - 监听进入目录时，自动检测 .xenv.toml 文件，并加载里面的配置
var InitDirenvCmd = &gcli.Command{
	Name:   "init-direnv",
	Desc:   "Initialize direnv state on current workdir",
	Hidden: true, // This is an internal command
	Func: func(c *gcli.Command, args []string) error {
		// Create tool service
		toolSvc, err := xenv.ToolService()
		if err != nil {
			return err
		}
		script, err1 := toolSvc.SetupDirenv()
		if err1 == nil {
			shell.OutputScript(script)
		}
		return err1
	},
}
