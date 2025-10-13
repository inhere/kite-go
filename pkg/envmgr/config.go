package envmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/gookit/goutil/fsutil"
)

// DefaultConfigManager 默认配置管理器实现
type DefaultConfigManager struct {
	configFile string
	baseDir    string
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configFile, baseDir string) *DefaultConfigManager {
	return &DefaultConfigManager{
		configFile: configFile,
		baseDir:    baseDir,
	}
}

// LoadConfig 加载配置
func (cm *DefaultConfigManager) LoadConfig() (*ShellEnvConfig, error) {
	// 解析路径中的变量
	configPath := cm.resolvePath(cm.configFile)

	if !fsutil.IsFile(configPath) {
		return cm.getDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config ShellEnvConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// 解析配置中的路径变量
	if err := cm.resolveConfigPaths(&config); err != nil {
		return nil, fmt.Errorf("failed to resolve config paths: %w", err)
	}

	// 验证配置
	if err := cm.ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置
func (cm *DefaultConfigManager) SaveConfig(config *ShellEnvConfig) error {
	if err := cm.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	configPath := cm.resolvePath(cm.configFile)

	// 确保目录存在
	if err := fsutil.Mkdir(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetSDKConfig 获取SDK配置
func (cm *DefaultConfigManager) GetSDKConfig(name string) (*SDKConfig, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return nil, err
	}

	for _, sdk := range config.SDKs {
		if sdk.Name == name {
			return &sdk, nil
		}
	}

	return nil, fmt.Errorf("SDK config not found: %s", name)
}

// ValidateConfig 验证配置
func (cm *DefaultConfigManager) ValidateConfig(config *ShellEnvConfig) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	// 验证SDK目录
	if config.SDKDir != "" {
		if !filepath.IsAbs(config.SDKDir) && !strings.HasPrefix(config.SDKDir, "~") {
			return fmt.Errorf("sdk_dir must be absolute path: %s", config.SDKDir)
		}
	}

	// 验证SDK配置
	sdkNames := make(map[string]bool)
	for i, sdk := range config.SDKs {
		if sdk.Name == "" {
			return fmt.Errorf("SDK name is empty at index %d", i)
		}

		if !IsValidSDKName(sdk.Name) {
			return fmt.Errorf("invalid SDK name: %s", sdk.Name)
		}

		if sdkNames[sdk.Name] {
			return fmt.Errorf("duplicate SDK name: %s", sdk.Name)
		}
		sdkNames[sdk.Name] = true

		if sdk.InstallDir == "" {
			return fmt.Errorf("SDK install_dir is empty for %s", sdk.Name)
		}
	}

	// 验证环境变量名
	for name := range config.AddEnvs {
		if !isValidEnvName(name) {
			return fmt.Errorf("invalid environment variable name: %s", name)
		}
	}

	for _, name := range config.RemoveEnvs {
		if !isValidEnvName(name) {
			return fmt.Errorf("invalid environment variable name: %s", name)
		}
	}

	return nil
}

// getDefaultConfig 获取默认配置
func (cm *DefaultConfigManager) getDefaultConfig() *ShellEnvConfig {
	return &ShellEnvConfig{
		AddPaths:   []string{},
		AddEnvs:    map[string]string{},
		RemoveEnvs: []string{},
		SDKDir:     fsutil.HomeDir() + "/.kite-go/sdk",
		SDKs:       []SDKConfig{},
	}
}

// resolveConfigPaths 解析配置中的路径变量
func (cm *DefaultConfigManager) resolveConfigPaths(config *ShellEnvConfig) error {
	// 解析SDK目录路径
	config.SDKDir = cm.resolvePath(config.SDKDir)

	// 解析添加的路径
	for i, path := range config.AddPaths {
		config.AddPaths[i] = cm.resolvePath(path)
	}

	// 解析SDK配置中的路径
	for i := range config.SDKs {
		config.SDKs[i].InstallDir = cm.resolvePath(config.SDKs[i].InstallDir)
	}

	return nil
}

// resolvePath 解析路径中的变量
func (cm *DefaultConfigManager) resolvePath(path string) string {
	if path == "" {
		return path
	}

	// 替换预定义变量
	path = strings.ReplaceAll(path, "$base", cm.baseDir)
	path = strings.ReplaceAll(path, "$config", filepath.Join(cm.baseDir, "config"))
	path = strings.ReplaceAll(path, "$data", filepath.Join(cm.baseDir, "data"))
	path = strings.ReplaceAll(path, "$tmp", filepath.Join(cm.baseDir, "tmp"))

	// 展开用户主目录
	if strings.HasPrefix(path, "~") {
		homeDir := fsutil.HomeDir()
		if path == "~" {
			path = homeDir
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(homeDir, path[2:])
		}
	}

	// 解析为绝对路径
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err == nil {
			path = abs
		}
	}

	return path
}

// isValidEnvName 检查环境变量名是否有效
func isValidEnvName(name string) bool {
	if name == "" {
		return false
	}

	// 环境变量名只能包含字母、数字和下划线，且不能以数字开头
	for i, r := range name {
		if i == 0 {
			if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_') {
				return false
			}
		} else {
			if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
				return false
			}
		}
	}

	return true
}

// GetSupportedSDKs 获取支持的SDK列表
func (cm *DefaultConfigManager) GetSupportedSDKs() ([]string, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return nil, err
	}

	var sdks []string
	for _, sdk := range config.SDKs {
		sdks = append(sdks, sdk.Name)
	}

	return sdks, nil
}

// AddSDKConfig 添加SDK配置
func (cm *DefaultConfigManager) AddSDKConfig(sdk SDKConfig) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	// 检查是否已存在
	for i, existing := range config.SDKs {
		if existing.Name == sdk.Name {
			config.SDKs[i] = sdk // 更新现有配置
			return cm.SaveConfig(config)
		}
	}

	// 添加新配置
	config.SDKs = append(config.SDKs, sdk)
	return cm.SaveConfig(config)
}

// RemoveSDKConfig 移除SDK配置
func (cm *DefaultConfigManager) RemoveSDKConfig(name string) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}

	for i, sdk := range config.SDKs {
		if sdk.Name == name {
			config.SDKs = append(config.SDKs[:i], config.SDKs[i+1:]...)
			return cm.SaveConfig(config)
		}
	}

	return fmt.Errorf("SDK config not found: %s", name)
}
