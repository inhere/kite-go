package models

import "time"

// Configuration 代表用户的配置信息，包含工具管理设置、路径配置、环境激活状态等数据
type Configuration struct {
	ID              string                  `json:"id"`
	BinDir          string                  `json:"bin_dir"`           // 默认: ~/.local/bin
	InstallDir      string                  `json:"install_dir"`       // 默认: ~/.xenv/tools
	ShellScriptsDir string                  `json:"shell_scripts_dir"` // 默认: ~/.config/xenv/hooks/
	Tools           []ToolChain             `json:"tools"`             // 可管理的工具链列表
	GlobalEnv       map[string]EnvironmentVariable `json:"global_env"` // 全局环境变量
	GlobalPaths     []PathEntry             `json:"global_paths"`      // 全局PATH条目
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}