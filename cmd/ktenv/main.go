package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/internal/cli/devcmd/envcmd"
	"github.com/inhere/kite-go/pkg/envmgr"
)

// main ktenv程序入口
func main() {
	em := envcmd.NewEnvManageCmd()
	em.Name = "ktenv"

	em.MustRun(nil)
}

// getBaseDir 获取基础目录
func getBaseDir() string {
	// 从环境变量获取
	if baseDir := os.Getenv("KITE_BASE_DIR"); baseDir != "" {
		return baseDir
	}

	// 默认目录
	return filepath.Join(fsutil.HomeDir(), ".kite-go")
}

// processCommand 处理命令
func processCommand(em *envmgr.DefaultEnvManager, cmd string, args []string) (string, error) {
	switch cmd {
	case "use":
		return processUseCommand(em, args)
	case "unuse":
		return processUnuseCommand(em, args)
	case "add":
		return processAddCommand(em, args)
	case "list":
		return processListCommand(em, args)
	case "help", "--help", "-h":
		return getHelpText(), nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

// processUseCommand 处理use命令
func processUseCommand(em *envmgr.DefaultEnvManager, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: ktenv use <sdk:version> [sdk:version...]")
	}

	var saveFlag bool
	var specs []string

	// 解析参数
	for _, arg := range args {
		if arg == "-s" || arg == "--save" {
			saveFlag = true
		} else {
			specs = append(specs, arg)
		}
	}

	if len(specs) == 0 {
		return "", fmt.Errorf("no SDK specifications provided")
	}

	// 解析版本规格
	versionSpecs, err := envmgr.ParseMultipleVersionSpecs(specs)
	if err != nil {
		return "", fmt.Errorf("invalid version specification: %w", err)
	}

	// 激活所有指定的SDK
	for _, spec := range versionSpecs {
		if err := em.UseSDK(spec.SDK, spec.Version, saveFlag); err != nil {
			return "", fmt.Errorf("failed to use %s: %w", spec.String(), err)
		}
	}

	// 生成当前shell的环境设置脚本
	shellType := envmgr.DetectShellType()
	return generateCurrentEnvScript(em, shellType)
}

// processUnuseCommand 处理unuse命令
func processUnuseCommand(em *envmgr.DefaultEnvManager, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: ktenv unuse <sdk> [sdk...]")
	}

	for _, sdk := range args {
		if err := em.UnuseSDK(sdk); err != nil {
			return "", fmt.Errorf("failed to unuse %s: %w", sdk, err)
		}
	}

	// 生成当前shell的环境设置脚本
	shellType := envmgr.DetectShellType()
	return generateCurrentEnvScript(em, shellType)
}

// processAddCommand 处理add命令
func processAddCommand(em *envmgr.DefaultEnvManager, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: ktenv add <sdk:version> [sdk:version...]")
	}

	// 解析版本规格
	versionSpecs, err := envmgr.ParseMultipleVersionSpecs(args)
	if err != nil {
		return "", fmt.Errorf("invalid version specification: %w", err)
	}

	// 下载安装所有指定的SDK
	for _, spec := range versionSpecs {
		if err := em.AddSDK(spec.SDK, spec.Version); err != nil {
			return "", fmt.Errorf("failed to add %s: %w", spec.String(), err)
		}
	}

	return fmt.Sprintf("Successfully installed %d SDKs\n", len(versionSpecs)), nil
}

// processListCommand 处理list命令
func processListCommand(em *envmgr.DefaultEnvManager, args []string) (string, error) {
	var sdkType string
	if len(args) > 0 {
		sdkType = args[0]
	}

	sdks, err := em.ListSDKs(sdkType)
	if err != nil {
		return "", fmt.Errorf("failed to list SDKs: %w", err)
	}

	if len(sdks) == 0 {
		return "No SDKs found\n", nil
	}

	var result string
	result += "Installed SDKs:\n"

	for _, sdk := range sdks {
		status := "inactive"
		if sdk.IsActive {
			status = "active"
		}

		if sdk.Installed {
			result += fmt.Sprintf("  %s:%s (%s) - %s\n", sdk.Name, sdk.Version, status, sdk.Path)
		} else {
			result += fmt.Sprintf("  %s (not installed)\n", sdk.Name)
		}
	}

	return result, nil
}

// generateCurrentEnvScript 生成当前环境脚本（只包含变更部分）
func generateCurrentEnvScript(em *envmgr.DefaultEnvManager, shellType envmgr.ShellType) (string, error) {
	state, err := em.GetActiveState()
	if err != nil {
		return "", err
	}

	shellGen := envmgr.NewShellGenerator("kite")
	var script string

	// 生成PATH更新
	if len(state.AddPaths) > 0 {
		pathScript, err := shellGen.GeneratePathUpdate(shellType, state.AddPaths, envmgr.PathAdd)
		if err != nil {
			return "", err
		}
		script += pathScript
	}

	// 生成环境变量设置
	if len(state.AddEnvs) > 0 {
		envScript, err := shellGen.GenerateEnvVars(shellType, state.AddEnvs)
		if err != nil {
			return "", err
		}
		script += envScript
	}

	return script, nil
}

// getHelpText 获取帮助文本
func getHelpText() string {
	return `ktenv - Kite Environment Manager

Usage:
  ktenv <command> [arguments]

Commands:
  use <sdk:version>...     Activate SDK versions
    -s, --save             Save configuration to project file
  unuse <sdk>...           Deactivate SDKs
  add <sdk:version>...     Download and install SDK versions
  list [sdk]               List installed SDKs
  help                     Show this help message

Examples:
  ktenv use node:18 go:1.21
  ktenv use -s node:lts
  ktenv unuse node
  ktenv add go:1.22
  ktenv list
  ktenv list go

Supported SDKs:
  go, node, java, flutter

Version formats:
  <sdk>:<version>         Exact version (go:1.21.5)
  <sdk>:<major>           Latest patch version (node:18)
  <sdk>:lts               Long-term support version
  <sdk>:latest            Latest stable version
`
}
