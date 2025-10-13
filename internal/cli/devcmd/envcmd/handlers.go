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
	str, err := em.ShowInstalledSDKs(envListOpts.sdkType)

	if err != nil {
		return fmt.Errorf("failed to list SDKs: %w", err)
	}
	c.Println(str)
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

var envConfigOpts = struct {
	edit bool
}{}

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
	show.JSON(state)
	return nil
}
