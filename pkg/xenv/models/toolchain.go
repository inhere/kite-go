package models

import "time"

// ToolChain 代表特定版本的开发工具（如Go、Node.js等），包含版本号、安装路径、别名等属性
type ToolChain struct {
	ID          string            `json:"id"`          // 唯一标识符，格式为 name:version
	Name        string            `json:"name"`        // 工具名称，如 "go", "node"
	Version     string            `json:"version"`     // 版本号，如 "1.21", "lts", "latest"
	Alias       []string          `json:"alias"`       // 别名列表，如 ["golang"] for go
	InstallURL  string            `json:"install_url"` // 可选，下载URL模板
	InstallDir  string            `json:"install_dir"` // 安装目录路径
	ActiveEnv   map[string]string `json:"active_env"`  // 激活时设置的额外环境变量
	Installed   bool              `json:"installed"`   // 是否已安装
	BinPaths    []string          `json:"bin_paths"`   // 该工具的二进制文件路径列表
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}