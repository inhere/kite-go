package models

import "time"

// LocalTools 代表本地已安装的工具链信息. 会保存到 ~/.xenv/tools/local.json
type LocalTools struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// SDK 工具链列表
	SdkTools    []InstalledTool `json:"sdk_tools"`
	SimpleTools []InstalledTool `json:"simple_tools"`
}

// FindSdkTool 根据名称和版本查询已安装的工具链
func (lt *LocalTools) FindSdkTool(name string, version string) *InstalledTool {
	for i, tool := range lt.SdkTools {
		if tool.Name == name && tool.Version == version {
			tool.Index = i
			return &tool
		}
	}
	return nil
}

// ListSdkTools 根据名称返回所有已安装的SDK工具链版本
func (lt *LocalTools) ListSdkTools(name string) []InstalledTool {
	var tools []InstalledTool
	for i, tool := range lt.SdkTools {
		if tool.Name == name {
			tool.Index = i
			tools = append(tools, tool)
		}
	}
	return tools
}

// FindSimpleTool 根据名称查询已安装的工具链
func (lt *LocalTools) FindSimpleTool(name string) *InstalledTool {
	for _, tool := range lt.SimpleTools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

// InstalledTool 代表已安装的工具链信息
type InstalledTool struct {
	ID         string    `json:"id"` // 唯一标识符，格式为 name:version
	Index int `json:"index"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	InstallDir string    `json:"install_dir"`
	BinDir     string    `json:"bin_dir"`
	BinPaths   []string  `json:"bin_paths"`
	Source string `json:"source"`
	Simple bool   `json:"simple"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
