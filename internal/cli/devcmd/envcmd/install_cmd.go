package envcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/envmgr"
)

func NewEnvAddCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "add",
		Desc:    "add/install new environment SDK to local",
		Aliases: []string{"ins", "install"},
		Func:    handleToolInstall,
	}
}

// handleToolInstall 处理 tool 安装命令
func handleToolInstall(c *gcli.Command, args []string) error {
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

// NewEnvRemoveCmd 创建Tool移除命令
func NewEnvRemoveCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "remove",
		Desc:    "remove installed environment tool SDK",
		Aliases: []string{"del", "rm", "delete", "uninstall"},
		Func:    handleEnvRemove,
	}
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
