package models

import (
	"time"

	"github.com/gookit/goutil/arrutil"
)

// ActivityState 代表用户当前激活的工具链和环境状态.
//  - 全局的会保存到 ~/.xenv/.xenv.toml
type ActivityState struct {
	ID string `json:"id"` // name or file path
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 激活的工具链映射 key为工具名，value为版本
	SDKs map[string]string `json:"active_sdks"`
	// 激活的环境变量
	Envs map[string]string `json:"envs"`
	// 激活的路径列表
	Paths []string `json:"paths"`
	// Tools 需要的工具列表
	Tools map[string]string `json:"tools"`
}

// NewActivityState creates a new ActivityState
func NewActivityState(id ...string) *ActivityState {
	return &ActivityState{
		ID:        arrutil.FirstOr(id, "default"),
		SDKs:      make(map[string]string),
		Envs:      make(map[string]string),
		Paths:     []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
