package models

import (
	"sort"
	"time"

	"github.com/inhere/kite-go/pkg/util"
)

const (
	InstalledMetaFile = "~/.xenv/tools.local.json"
)

// ToolsLocal 代表本地已安装的工具链信息. 会保存到 ~/.xenv/tools.local.json
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
	ID         string `json:"id"` // 唯一标识符，格式为 name:version
	Name       string `json:"name"`
	Version    string `json:"version"`
	// InstallDir 当前版本的工具安装目录路径
	InstallDir string `json:"install_dir"`
	// BinDir 可执行文件目录, 相对于 InstallDir
	//  - 为空时，默认为 install_dir/bin
	BinDir    string    `json:"bin_dir,omitempty"`
	Source    string    `json:"source"`
	IsSDK     bool      `json:"is_sdk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 内部使用
	Index  int        `json:"-"` // 表示在列表中的索引
	Config *ToolChain `json:"-"` // 工具链配置
}

// BinDirPath 返回可执行文件目录的绝对路径
func (t *InstalledTool) BinDirPath() string {
	return util.NormalizePath(t.binDirPath())
}

// 返回可执行文件目录的绝对路径
func (t *InstalledTool) binDirPath() string {
	if t.BinDir == "" {
		return t.InstallDir + "/bin"
	}
	return t.InstallDir + "/" + t.BinDir
}

// RenderActiveEnv 渲染工具链的激活环境变量
func (t *InstalledTool) RenderActiveEnv() map[string]string {
	if len(t.Config.ActiveEnv) == 0 {
		return nil
	}

	varMap := map[string]string{
		"name":        t.Name,
		"version":     t.Version,
		"install_dir": util.NormalizePath(t.InstallDir),
	}
	return t.Config.RenderActiveEnv(varMap)
}
