package models

import (
	"github.com/gookit/goutil/arrutil"
)

const (
	// GlobalStateFile global state file path
	GlobalStateFile = "~/.config/xenv/global.toml"
	// LocalStateFile local state file path
	LocalStateFile = ".xenv.toml"
)

// ActivityState 代表用户当前激活的工具链和环境状态.
//  - 全局的会保存到 ~/.config/xenv/global.toml
type ActivityState struct {
	// 激活的 ENV 路径列表
	Paths []string `json:"paths"`
	// 激活的工具链映射 key为工具名，value为版本
	SDKs map[string]string `json:"sdks"`
	// 激活的环境变量
	Envs map[string]string `json:"envs"`
	// Tools 需要的工具列表
	Tools map[string]string `json:"tools"`
	// enable_global 是否启用全局环境配置
	// EnableGlobal bool `json:"enable_global"`
	File string `json:"-"` // state file path
	// 本次改变的数据,保存后置为nil
	ChangeData *ActivityState `json:"-"`
	// 创建时间
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

// NewActivityState creates a new ActivityState
func NewActivityState(filePath ...string) *ActivityState {
	return &ActivityState{
		File:  arrutil.FirstOr(filePath),
		SDKs:  make(map[string]string),
		Envs:  make(map[string]string),
		Paths: []string{},
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}
}

// AddToolsWithEnvsPaths 新增激活工具和配置ENV, PATH
func (as *ActivityState) AddToolsWithEnvsPaths(tools, envs map[string]string, paths []string) {
	for name, version := range tools {
		as.SDKs[name] = version
	}
	for name, value := range envs {
		as.Envs[name] = value
	}
	for _, path := range paths {
		as.AddActivePath(path)
	}
}

// Merge other to current. 合并两个状态数据
func (as *ActivityState) Merge(other *ActivityState) {
	if other == nil {
		return
	}
	as.AddToolsWithEnvsPaths(other.SDKs, other.Envs, other.Paths)
}

// DelToolsWithEnvsPaths 删除激活工具和相关的 ENV, PATH
func (as *ActivityState) DelToolsWithEnvsPaths(toolNames, envNames, paths []string) {
	if len(toolNames) > 0 {
		for _, name := range toolNames {
			delete(as.SDKs, name)
		}
	}

	if len(envNames) > 0 {
		for _, name := range envNames {
			delete(as.Envs, name)
		}
	}

	as.RemovePaths(paths)
}

// AddActivePath 添加激活路径, 会先检测是否已存在
func (as *ActivityState) AddActivePath(path string) {
	// 检查路径是否已存在
	for _, p := range as.Paths {
		if p == path {
			return
		}
	}
	as.Paths = append(as.Paths, path)
}

// RemovePaths 删除激活路径
func (as *ActivityState) RemovePaths(paths []string) {
	if len(paths) == 0 {
		return
	}

	var newPaths []string
	for _, path := range as.Paths {
		if arrutil.StringsContains(paths, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	as.Paths = newPaths
}

// RemoveTool 删除激活的工具
func (as *ActivityState) RemoveTool(name string) bool {
	_, exists := as.SDKs[name]
	if exists {
		delete(as.SDKs, name)
	}
	return exists
}
