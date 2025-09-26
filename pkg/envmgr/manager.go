package envmgr

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultEnvManager 默认环境管理器实现
type DefaultEnvManager struct {
	configManager ConfigManager
	stateManager  StateManager
	sdkManager    SDKManager
	shellGen      ShellGenerator
	bashDir       string
	kiteCommand   string
}

// NewEnvManager 创建环境管理器
func NewEnvManager(baseDir, kiteCommand string) *DefaultEnvManager {
	// 配置文件路径
	configFile := filepath.Join(baseDir, "config", "module", "shell_env.yml")
	stateFile := filepath.Join(baseDir, "data", "shell_env", "active.json")
	fmt.Println("configFile: %s, stateFile: %s", configFile, stateFile)

	// 创建管理器实例
	configManager := NewConfigManager(configFile, baseDir)
	stateManager := NewStateManager(stateFile)
	sdkManager := NewSDKManager(configManager)
	shellGen := NewShellGenerator(kiteCommand)

	return &DefaultEnvManager{
		configManager: configManager,
		stateManager:  stateManager,
		sdkManager:    sdkManager,
		shellGen:      shellGen,
		bashDir:       baseDir,
		kiteCommand:   kiteCommand,
	}
}

// UseSDK 激活指定SDK版本
func (em *DefaultEnvManager) UseSDK(sdk, version string, save bool) error {
	// 验证SDK名称和版本
	if !IsValidSDKName(sdk) {
		return fmt.Errorf("invalid SDK name: %s", sdk)
	}

	if !IsValidVersion(version) {
		return fmt.Errorf("invalid version: %s", version)
	}

	// 检查SDK是否支持
	_, err := em.configManager.GetSDKConfig(sdk)
	if err != nil {
		return fmt.Errorf("unsupported SDK: %s", sdk)
	}

	// 检查SDK是否已安装
	if !em.sdkManager.IsInstalled(sdk, version) {
		return fmt.Errorf("SDK %s:%s is not installed, use 'ktenv add %s:%s' to install it first", sdk, version, sdk, version)
	}

	// 验证SDK安装
	if err := em.sdkManager.ValidateSDK(sdk, version); err != nil {
		return fmt.Errorf("SDK %s:%s validation failed: %w", sdk, version, err)
	}

	// 更新活跃状态
	if err := em.stateManager.UpdateSDKState(sdk, version, true); err != nil {
		return fmt.Errorf("failed to update SDK state: %w", err)
	}

	// 获取SDK路径并添加到PATH
	binPath := em.sdkManager.GetSDKBinPath(sdk, version)
	if binPath != "" {
		if err := em.stateManager.AddPath(binPath); err != nil {
			return fmt.Errorf("failed to add SDK bin path: %w", err)
		}
	}

	// 设置SDK特定的环境变量
	envVars, err := em.sdkManager.GetSDKEnvVars(sdk, version)
	if err != nil {
		return fmt.Errorf("failed to get SDK env vars: %w", err)
	}

	for name, value := range envVars {
		if err := em.stateManager.SetEnv(name, value); err != nil {
			return fmt.Errorf("failed to set env var %s: %w", name, err)
		}
	}

	// 如果需要保存到配置文件
	if save {
		if err := em.saveCurrentState(); err != nil {
			return fmt.Errorf("failed to save state to config: %w", err)
		}
	}

	return nil
}

// UnuseSDK 取消激活SDK
func (em *DefaultEnvManager) UnuseSDK(sdk string) error {
	// 检查SDK是否处于激活状态
	currentSDKs, err := em.stateManager.GetCurrentSDKs()
	if err != nil {
		return fmt.Errorf("failed to get current SDKs: %w", err)
	}

	version, isActive := currentSDKs[sdk]
	if !isActive {
		return fmt.Errorf("SDK %s is not currently active", sdk)
	}

	// 移除SDK状态
	if err := em.stateManager.UpdateSDKState(sdk, "", false); err != nil {
		return fmt.Errorf("failed to update SDK state: %w", err)
	}

	// 移除SDK的bin路径
	binPath := em.sdkManager.GetSDKBinPath(sdk, version)
	if binPath != "" {
		if err := em.stateManager.RemovePath(binPath); err != nil {
			return fmt.Errorf("failed to remove SDK bin path: %w", err)
		}
	}

	// 移除SDK特定的环境变量
	envVars, err := em.sdkManager.GetSDKEnvVars(sdk, version)
	if err != nil {
		return fmt.Errorf("failed to get SDK env vars: %w", err)
	}

	for name := range envVars {
		if err := em.stateManager.UnsetEnv(name); err != nil {
			return fmt.Errorf("failed to unset env var %s: %w", name, err)
		}
	}

	return nil
}

// AddSDK 下载安装SDK
func (em *DefaultEnvManager) AddSDK(sdk, version string) error {
	// 验证SDK名称和版本
	if !IsValidSDKName(sdk) {
		return fmt.Errorf("invalid SDK name: %s", sdk)
	}

	if !IsValidVersion(version) {
		return fmt.Errorf("invalid version: %s", version)
	}

	// 检查SDK是否支持
	_, err := em.configManager.GetSDKConfig(sdk)
	if err != nil {
		return fmt.Errorf("unsupported SDK: %s", sdk)
	}

	// 检查SDK是否已安装
	if em.sdkManager.IsInstalled(sdk, version) {
		return fmt.Errorf("SDK %s:%s is already installed", sdk, version)
	}

	// 下载并安装SDK
	if err := em.sdkManager.InstallSDK(sdk, version); err != nil {
		return fmt.Errorf("failed to install SDK %s:%s: %w", sdk, version, err)
	}

	// 验证安装
	if err := em.sdkManager.ValidateSDK(sdk, version); err != nil {
		return fmt.Errorf("SDK %s:%s installation validation failed: %w", sdk, version, err)
	}

	return nil
}

// ListSDKs 列出已安装的SDK
func (em *DefaultEnvManager) ListSDKs(sdkType string) ([]SDKInfo, error) {
	var result []SDKInfo

	// 获取当前激活的SDK
	currentSDKs, err := em.stateManager.GetCurrentSDKs()
	if err != nil {
		return nil, fmt.Errorf("failed to get current SDKs: %w", err)
	}

	// 获取支持的SDK列表
	supportedSDKs, err := em.configManager.GetSupportedSDKs()
	if err != nil {
		return nil, fmt.Errorf("failed to get supported SDKs: %w", err)
	}

	for _, sdk := range supportedSDKs {
		// 如果指定了SDK类型，只列出匹配的SDK
		if sdkType != "" && !strings.EqualFold(sdk, sdkType) {
			continue
		}

		// 获取已安装的版本
		versions, err := em.sdkManager.ListVersions(sdk)
		if err != nil {
			continue // 忽略错误，继续处理其他SDK
		}

		for _, version := range versions {
			isActive := currentSDKs[sdk] == version
			sdkInfo := SDKInfo{
				Name:      sdk,
				Version:   version,
				IsActive:  isActive,
				Path:      em.sdkManager.GetSDKPath(sdk, version),
				Installed: em.sdkManager.IsInstalled(sdk, version),
			}
			result = append(result, sdkInfo)
		}

		// 如果没有安装的版本，也显示SDK信息
		if len(versions) == 0 {
			sdkInfo := SDKInfo{
				Name:      sdk,
				Version:   "",
				IsActive:  false,
				Path:      "",
				Installed: false,
			}
			result = append(result, sdkInfo)
		}
	}

	return result, nil
}

// GetActiveState 获取当前活跃状态
func (em *DefaultEnvManager) GetActiveState() (*ActiveState, error) {
	return em.stateManager.LoadState()
}

// GenerateShellScript 生成shell脚本
func (em *DefaultEnvManager) GenerateShellScript(shellType ShellType) (string, error) {
	// 获取当前状态
	state, err := em.stateManager.LoadState()
	if err != nil {
		return "", fmt.Errorf("failed to load state: %w", err)
	}

	// 获取配置
	config, err := em.configManager.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// 合并配置和状态中的路径和环境变量
	mergedState := em.mergeConfigAndState(config, state)

	// 生成脚本
	script, err := em.shellGen.GenerateScript(shellType, mergedState, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate shell script: %w", err)
	}

	return script, nil
}

// ProcessKtenvCommand 处理ktenv命令
func (em *DefaultEnvManager) ProcessKtenvCommand(cmd string, args []string) (string, error) {
	switch cmd {
	case "use":
		return em.processUseCommand(args)
	case "unuse":
		return em.processUnuseCommand(args)
	case "add":
		return em.processAddCommand(args)
	case "list":
		return em.processListCommand(args)
	default:
		return "", fmt.Errorf("unknown ktenv command: %s", cmd)
	}
}

// processUseCommand 处理use命令
func (em *DefaultEnvManager) processUseCommand(args []string) (string, error) {
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
	versionSpecs, err := ParseMultipleVersionSpecs(specs)
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
	shellType := DetectShellType()
	script, err := em.generateCurrentEnvScript(shellType)
	if err != nil {
		return "", fmt.Errorf("failed to generate env script: %w", err)
	}

	return script, nil
}

// processUnuseCommand 处理unuse命令
func (em *DefaultEnvManager) processUnuseCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: ktenv unuse <sdk> [sdk...]")
	}

	for _, sdk := range args {
		if err := em.UnuseSDK(sdk); err != nil {
			return "", fmt.Errorf("failed to unuse %s: %w", sdk, err)
		}
	}

	// 生成当前shell的环境设置脚本
	shellType := DetectShellType()
	script, err := em.generateCurrentEnvScript(shellType)
	if err != nil {
		return "", fmt.Errorf("failed to generate env script: %w", err)
	}

	return script, nil
}

// processAddCommand 处理add命令
func (em *DefaultEnvManager) processAddCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: ktenv add <sdk:version> [sdk:version...]")
	}

	// 解析版本规格
	versionSpecs, err := ParseMultipleVersionSpecs(args)
	if err != nil {
		return "", fmt.Errorf("invalid version specification: %w", err)
	}

	// 下载安装所有指定的SDK
	for _, spec := range versionSpecs {
		if err := em.AddSDK(spec.SDK, spec.Version); err != nil {
			return "", fmt.Errorf("failed to add %s: %w", spec.String(), err)
		}
	}

	return fmt.Sprintf("Successfully installed %d SDKs", len(versionSpecs)), nil
}

// processListCommand 处理list命令
func (em *DefaultEnvManager) processListCommand(args []string) (string, error) {
	var sdkType string
	if len(args) > 0 {
		sdkType = args[0]
	}

	sdks, err := em.ListSDKs(sdkType)
	if err != nil {
		return "", fmt.Errorf("failed to list SDKs: %w", err)
	}

	if len(sdks) == 0 {
		return "No SDKs found", nil
	}

	var result strings.Builder
	result.WriteString("Installed SDKs:\n")

	for _, sdk := range sdks {
		status := "inactive"
		if sdk.IsActive {
			status = "active"
		}

		if sdk.Installed {
			result.WriteString(fmt.Sprintf("  %s:%s (%s) - %s\n", sdk.Name, sdk.Version, status, sdk.Path))
		} else {
			result.WriteString(fmt.Sprintf("  %s (not installed)\n", sdk.Name))
		}
	}

	return result.String(), nil
}

// mergeConfigAndState 合并配置和状态
func (em *DefaultEnvManager) mergeConfigAndState(config *ShellEnvConfig, state *ActiveState) *ActiveState {
	merged := state.Clone()

	// 合并配置中的路径
	pathSet := make(map[string]bool)
	for _, path := range merged.AddPaths {
		pathSet[path] = true
	}

	for _, path := range config.AddPaths {
		if !pathSet[path] {
			merged.AddPaths = append(merged.AddPaths, path)
			pathSet[path] = true
		}
	}

	// 合并配置中的环境变量
	for name, value := range config.AddEnvs {
		merged.AddEnvs[name] = value
	}

	return merged
}

// generateCurrentEnvScript 生成当前环境脚本（只包含变更部分）
func (em *DefaultEnvManager) generateCurrentEnvScript(shellType ShellType) (string, error) {
	state, err := em.stateManager.LoadState()
	if err != nil {
		return "", err
	}

	var script strings.Builder

	// 生成PATH更新
	if len(state.AddPaths) > 0 {
		pathScript, err := em.shellGen.GeneratePathUpdate(shellType, state.AddPaths, PathAdd)
		if err != nil {
			return "", err
		}
		script.WriteString(pathScript)
	}

	// 生成环境变量设置
	if len(state.AddEnvs) > 0 {
		envScript, err := em.shellGen.GenerateEnvVars(shellType, state.AddEnvs)
		if err != nil {
			return "", err
		}
		script.WriteString(envScript)
	}

	return script.String(), nil
}

// saveCurrentState 保存当前状态到配置文件
func (em *DefaultEnvManager) saveCurrentState() error {
	// 这里可以实现将当前活跃状态保存到项目配置文件的逻辑
	// 例如 .kite.yml 或 .envrc 文件
	return nil
}

// GetStats 获取环境管理统计信息
func (em *DefaultEnvManager) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 状态统计
	stateStats, err := em.stateManager.GetStateStats()
	if err != nil {
		return nil, err
	}
	stats["state"] = stateStats

	// SDK统计
	sdks, err := em.ListSDKs("")
	if err != nil {
		return nil, err
	}

	installedCount := 0
	activeCount := 0
	for _, sdk := range sdks {
		if sdk.Installed {
			installedCount++
		}
		if sdk.IsActive {
			activeCount++
		}
	}

	stats["sdks"] = map[string]interface{}{
		"total":     len(sdks),
		"installed": installedCount,
		"active":    activeCount,
	}

	return stats, nil
}
