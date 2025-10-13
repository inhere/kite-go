package envmgr

import (
	"time"
)

// ManagerConfig shell env manager configuration
type ManagerConfig struct {
	// AddPaths 添加到 PATH 的路径列表
	AddPaths []string `yaml:"add_paths" json:"add_paths"`
	// AddEnvs 添加的环境变量列表
	AddEnvs []string `yaml:"add_envs" json:"add_envs"`
	// ConfigFile 配置SDK目录，安装URL等配置YAML文件路径
	ConfigFile string `yaml:"config_file" json:"config_file"`
	// ActiveFile 当前环境的设置信息JSON文件路径
	ActiveFile string `yaml:"active_file" json:"active_file"`
	// CustomDir 自定义eval脚本目录（bash, pwsh等使用）
	CustomDir string `yaml:"custom_dir" json:"custom_dir"`
}

// SDKConfig SDK配置定义
type SDKConfig struct {
	Name       string            `yaml:"name" json:"name"`                             // SDK名称标识符
	InstallURL string            `yaml:"install_url,omitempty" json:"install_url"`     // 下载URL模板
	InstallDir string            `yaml:"install_dir" json:"install_dir"`               // 安装目录模板
	ActiveEnv  map[string]string `yaml:"active_env,omitempty" json:"active_env"`       // 激活时的额外环境变量
}

// ShellEnvConfig shell环境配置
type ShellEnvConfig struct {
	AddPaths   []string            `yaml:"add_paths" json:"add_paths"`     // 添加到PATH的路径列表
	AddEnvs    map[string]string   `yaml:"add_envs" json:"add_envs"`       // 需要设置的环境变量映射
	RemoveEnvs []string            `yaml:"remove_envs" json:"remove_envs"` // 需要移除的环境变量列表
	SDKDir     string              `yaml:"sdk_dir" json:"sdk_dir"`         // SDK基础目录
	SDKs       []SDKConfig         `yaml:"sdks" json:"sdks"`               // SDK配置列表
}

// ActiveState 活跃状态数据
type ActiveState struct {
	CurrentSDKs map[string]string `json:"current_sdks"` // 当前激活的SDK映射，key为SDK名，value为版本
	AddPaths    []string          `json:"add_paths"`    // 需要添加到PATH的路径列表
	AddEnvs     map[string]string `json:"add_envs"`     // 需要设置的环境变量映射
	UpdatedAt   time.Time         `json:"updated_at"`   // 最后更新时间
}

// ShellType shell类型枚举
type ShellType string

const (
	ShellBash       ShellType = "bash"
	ShellZsh        ShellType = "zsh"
	ShellPowerShell ShellType = "pwsh"
	ShellCmd        ShellType = "cmd"
)

// IsValid 检查shell类型是否有效
func (st ShellType) IsValid() bool {
	switch st {
	case ShellBash, ShellZsh, ShellPowerShell, ShellCmd:
		return true
	default:
		return false
	}
}

// String 返回shell类型字符串
func (st ShellType) String() string {
	return string(st)
}

// EnvManager 环境管理器接口
type EnvManager interface {
	// UseSDK 激活指定SDK版本
	UseSDK(sdk, version string, save bool) error

	// UnuseSDK 取消激活SDK
	UnuseSDK(sdk string) error

	// AddSDK 下载安装SDK
	AddSDK(sdk, version string) error

	// ListSDKs 列出已安装的SDK
	ListSDKs(sdkType string) ([]SDKInfo, error)

	// GetActiveState 获取当前活跃状态
	GetActiveState() (*ActiveState, error)

	// GenerateShellScript 生成shell脚本
	GenerateShellScript(shellType ShellType) (string, error)
}

// SDKInfo SDK信息
type SDKInfo struct {
	Name      string `json:"name"`      // SDK名称
	Version   string `json:"version"`   // 版本号
	IsActive  bool   `json:"is_active"` // 是否激活
	Path      string `json:"path"`      // 安装路径
	Installed bool   `json:"installed"` // 是否已安装
}

// VersionSpec 版本规格
type VersionSpec struct {
	SDK     string // SDK名称
	Version string // 版本规格
}

// ParseVersionSpec 解析版本规格 "sdk:version"
// 实现在 version.go 文件中

// SDKManager SDK管理器接口
type SDKManager interface {
	// DownloadSDK 下载SDK
	DownloadSDK(sdk, version string) error

	// InstallSDK 安装SDK
	InstallSDK(sdk, version string) error

	// UninstallSDK 卸载SDK
	UninstallSDK(sdk, version string) error

	// GetSDKPath 获取SDK路径
	GetSDKPath(sdk, version string) string

	// GetSDKBinPath 获取SDK的可执行文件路径
	GetSDKBinPath(sdk, version string) string

	// IsInstalled 检查SDK是否已安装
	IsInstalled(sdk, version string) bool

	// ListVersionDirs 列出本地的可用SDK版本 {version: dir_path}
	ListVersionDirs(sdk string) (map[string]string, error)

	// ListVersions 列出本地的可用SDK版本
	ListVersions(sdk string) ([]string, error)

	// ValidateSDK 验证SDK安装
	ValidateSDK(sdk, version string) error

	// GetSDKEnvVars 获取SDK需要设置的环境变量
	GetSDKEnvVars(sdk, version string) (map[string]string, error)
}

// ConfigManager 配置管理器接口
type ConfigManager interface {
	// LoadConfig 加载配置
	LoadConfig() (*ShellEnvConfig, error)

	// SaveConfig 保存配置
	SaveConfig(config *ShellEnvConfig) error

	// GetSDKConfig 获取SDK配置
	GetSDKConfig(name string) (*SDKConfig, error)

	// ValidateConfig 验证配置
	ValidateConfig(config *ShellEnvConfig) error
	// GetSupportedSDKs 获取支持的SDK列表
	GetSupportedSDKs() ([]string, error)
}

// StateManager 激活的sdk环境状态管理器接口
type StateManager interface {
	// LoadState 加载活跃状态
	LoadState() (*ActiveState, error)

	// SaveState 保存活跃状态
	SaveState(state *ActiveState) error

	// UpdateSDKState 更新SDK状态
	UpdateSDKState(sdk, version string, active bool) error

	// GetCurrentSDKs 获取当前激活的SDK
	GetCurrentSDKs() (map[string]string, error)

	// AddPath 添加路径到状态
	AddPath(path string) error
	// SetEnv 设置环境变量到状态
	SetEnv(name, value string) error
	// RemovePath 从状态中移除路径
	RemovePath(path string) error
	// UnsetEnv 从状态中移除环境变量
	UnsetEnv(name string) error
	// GetStateStats 获取状态统计信息
	GetStateStats() (*StateStats, error)
}

// ShellGenerator Shell脚本生成器接口
type ShellGenerator interface {
	// GenerateScript 生成shell脚本
	GenerateScript(shellType ShellType, state *ActiveState, config *ShellEnvConfig) (string, error)

	// GenerateKtenvFunction 生成ktenv函数
	GenerateKtenvFunction(shellType ShellType) (string, error)

	// GenerateEnvVars 生成环境变量设置
	GenerateEnvVars(shellType ShellType, envs map[string]string) (string, error)

	// GeneratePathUpdate 生成PATH更新脚本
	GeneratePathUpdate(shellType ShellType, paths []string, operation PathOperation) (string, error)
}

// PathOperation 路径操作类型
type PathOperation int

const (
	PathAdd    PathOperation = iota // 添加路径
	PathRemove                      // 移除路径
	PathSet                         // 设置路径
)
