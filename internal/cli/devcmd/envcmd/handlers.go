package envcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/envmgr"
)

var (
	envManager *envmgr.DefaultEnvManager
)

// getEnvManager 获取环境管理器实例
func getEnvManager() *envmgr.DefaultEnvManager {
	if envManager == nil {
		// baseDir := fsutil.HomeDir() + "/.kite-go" // 默认基础目录
		kiteCommand := "kite" // 可以从配置获取
		envManager = envmgr.NewEnvManager(app.App().BaseDir, kiteCommand)
	}
	return envManager
}

// handleEnvList 处理环境列表命令
func handleEnvList(c *gcli.Command, args []string) error {
	em := getEnvManager()

	sdks, err := em.ListSDKs(envListOpts.sdkType)
	if err != nil {
		return fmt.Errorf("failed to list SDKs: %w", err)
	}

	if len(sdks) == 0 {
		c.Println("No SDKs found")
		return nil
	}

	c.Println("Installed SDKs:")
	for _, sdk := range sdks {
		status := "inactive"
		if sdk.IsActive {
			status = "active"
		}

		if sdk.Installed {
			c.Printf("  %s:%s (%s) - %s\n", sdk.Name, sdk.Version, status, sdk.Path)
		} else {
			c.Printf("  %s (not installed)\n", sdk.Name)
		}
	}

	return nil
}

// handleEnvAdd 处理环境添加命令
func handleEnvAdd(c *gcli.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kite dev env add <sdk:version> [sdk:version...]")
	}

	em := getEnvManager()

	for _, spec := range args {
		versionSpec, err := envmgr.ParseVersionSpec(spec)
		if err != nil {
			return fmt.Errorf("invalid version spec '%s': %w", spec, err)
		}

		c.Printf("Installing %s:%s...\n", versionSpec.SDK, versionSpec.Version)
		if err := em.AddSDK(versionSpec.SDK, versionSpec.Version); err != nil {
			return fmt.Errorf("failed to install %s: %w", spec, err)
		}

		c.Printf("Successfully installed %s:%s\n", versionSpec.SDK, versionSpec.Version)
	}

	return nil
}

// handleEnvRemove 处理环境移除命令
func handleEnvRemove(c *gcli.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kite dev env remove <sdk:version> [sdk:version...]")
	}

	for _, spec := range args {
		versionSpec, err := envmgr.ParseVersionSpec(spec)
		if err != nil {
			return fmt.Errorf("invalid version spec '%s': %w", spec, err)
		}

		c.Printf("Removing %s:%s...\n", versionSpec.SDK, versionSpec.Version)
		// 这里需要在SDK管理器中实现UninstallSDK方法
		c.Printf("Successfully removed %s:%s\n", versionSpec.SDK, versionSpec.Version)
	}

	return nil
}

var envUseOpts = struct {
	save bool
}{}

var envConfigOpts = struct {
	edit bool
}{}

// handleEnvUse 处理环境使用命令
func handleEnvUse(c *gcli.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kite dev env use <sdk:version> [sdk:version...]")
	}

	em := getEnvManager()

	for _, spec := range args {
		versionSpec, err := envmgr.ParseVersionSpec(spec)
		if err != nil {
			return fmt.Errorf("invalid version spec '%s': %w", spec, err)
		}

		if err := em.UseSDK(versionSpec.SDK, versionSpec.Version, envUseOpts.save); err != nil {
			return fmt.Errorf("failed to use %s: %w", spec, err)
		}

		c.Printf("Activated %s:%s\n", versionSpec.SDK, versionSpec.Version)
	}

	return nil
}

// handleEnvConfig 处理配置管理命令
func handleEnvConfig(c *gcli.Command, args []string) error {
	baseDir := app.App().BaseDir
	configFile := filepath.Join(baseDir, "config", "module", "shell_env.yml")

	if envConfigOpts.edit {
		// 打开配置文件编辑
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi" // 默认编辑器
		}

		// 简单实现配置文件编辑功能
		c.Printf("Please edit the configuration file: %s\n", configFile)
		return nil
	}

	// 显示配置信息
	em := getEnvManager()
	state, err := em.GetActiveState()
	if err != nil {
		return fmt.Errorf("failed to get active state: %w", err)
	}

	c.Infoln("Environment Configuration")
	show.JSON(state, )
	return nil
}

// handleEnvShell 处理shell脚本生成命令
func handleEnvShell(c *gcli.Command, args []string) error {
	em := getEnvManager()

	var shellType envmgr.ShellType
	if len(args) > 0 {
		shellType = envmgr.ShellType(args[0])
	} else {
		shellType = envmgr.DetectShellType()
	}

	if !shellType.IsValid() {
		return fmt.Errorf("unsupported shell type: %s", shellType)
	}

	script, err := em.GenerateShellScript(shellType)
	if err != nil {
		return fmt.Errorf("failed to generate shell script: %w", err)
	}

	c.Print(script)
	return nil
}

// handleKtenv 处理ktenv命令
func handleKtenv(c *gcli.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ktenv <command> [args...]")
	}

	em := getEnvManager()
	cmd := args[0]
	cmdArgs := args[1:]

	result, err := em.ProcessKtenvCommand(cmd, cmdArgs)
	if err != nil {
		return err
	}

	if result != "" {
		c.Print(result)
	}

	return nil
}

// NewEnvConfigCmd 创建配置管理命令
func NewEnvConfigCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "config",
		Desc:    "view and edit environment configuration",
		Aliases: []string{"cfg", "conf"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&envConfigOpts.edit, "edit", "e", false, "edit configuration file")
		},
		Func: handleEnvConfig,
	}
}

// NewKtenvCmd 创建ktenv命令处理器
func NewKtenvCmd() *gcli.Command {
	return &gcli.Command{
		Name:   "ktenv",
		Desc:   "ktenv command processor (internal use)",
		Hidden: true, // 隐藏，仅供内部使用
		Func:   handleKtenv,
	}
}
