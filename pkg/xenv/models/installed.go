package models

import "time"

// LocalTools 代表本地已安装的工具链信息. 会保存到 ~/.xenv/tools/local.json
type LocalTools struct {
	Tools []InstalledTool `json:"tools"`
	SimpleTools []InstalledTool `json:"simple_tools"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InstalledTool 代表已安装的工具链信息
type InstalledTool struct {
	ID         string    `json:"id"` // 唯一标识符，格式为 name:version
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	InstallDir string    `json:"install_dir"`
	BinDir     string    `json:"bin_dir"`
	BinPaths   []string  `json:"bin_paths"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
