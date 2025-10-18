package models

import (
	"sort"
	"time"
)

// ToolsLocal 代表本地已安装的工具链信息. 会保存到 ~/.xenv/tools/local.json
type ToolsLocal struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// SDKs SDK工具链列表. 允许多版本存在
	SDKs []InstalledTool `json:"sdks"`
	// Tools 简单工具列表. 只允许一个版本存在
	Tools []InstalledTool `json:"tools"`
}

// FindSdkByID 根据ID(=名称:版本)查询已安装的工具链
func (lt *ToolsLocal) FindSdkByID(toolId string) *InstalledTool {
	for i, tool := range lt.SDKs {
		if tool.ID == toolId {
			tool.Index = i
			return &tool
		}
	}
	return nil
}

// FindSdkByNameAndVersion 根据名称和版本查询已安装的工具链
func (lt *ToolsLocal) FindSdkByNameAndVersion(name string, version string) *InstalledTool {
	for i, tool := range lt.SDKs {
		if tool.Name == name && tool.Version == version {
			tool.Index = i
			return &tool
		}
	}
	return nil
}

// ListSdkByName 根据名称返回所有已安装的SDK工具链版本
func (lt *ToolsLocal) ListSdkByName(name string) []InstalledTool {
	var tools []InstalledTool
	for i, tool := range lt.SDKs {
		if tool.Name == name {
			tool.Index = i
			tools = append(tools, tool)
		}
	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Version > tools[j].Version
	})
	return tools
}

// FindToolByName 根据名称查询已安装的简单工具
func (lt *ToolsLocal) FindToolByName(name string) *InstalledTool {
	for _, tool := range lt.Tools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

// InstalledTool 代表已安装的工具链信息
type InstalledTool struct {
	ID         string    `json:"id"` // 唯一标识符，格式为 name:version
	Index    int      `json:"-"`     // 内部使用，表示在列表中的索引
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	InstallDir string    `json:"install_dir"`
	BinDir   string   `json:"bin_dir,omitempty"`
	BinPaths []string `json:"bin_paths,omitempty"`
	Source   string   `json:"source"`
	IsSDK    bool     `json:"is_sdk"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
