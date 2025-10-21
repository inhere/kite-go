package models

import (
	"time"

	"github.com/gookit/goutil/arrutil"
)

// ActivityState 代表用户当前激活的工具链和环境状态. 全局的会保存到 ~/.config/xenv/activity.json
type ActivityState struct {
	ID string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 激活的工具链映射 key为工具名，value为版本
	ActiveTools map[string]string `json:"active_tools"`
	ActiveEnv   map[string]string `json:"active_env"`   // 激活的环境变量
	ActivePaths []string          `json:"active_paths"` // 激活的路径列表
}

func NewActivityState() *ActivityState {
	return &ActivityState{
		ID:          "default",
		ActiveTools: make(map[string]string),
		ActiveEnv:   make(map[string]string),
		ActivePaths: []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// AddToolsWithEnvsPaths 新增激活工具和配置ENV, PATH
func (as *ActivityState) AddToolsWithEnvsPaths(tools, envs map[string]string, paths []string) {
	for name, version := range tools {
		as.ActiveTools[name] = version
	}
	for name, value := range envs {
		as.ActiveEnv[name] = value
	}
	for _, path := range paths {
		as.AddActivePath(path)
	}
}

// DelToolsWithEnvsPaths 删除激活工具和相关的 ENV, PATH
func (as *ActivityState) DelToolsWithEnvsPaths(toolNames, envNames, paths []string) {
	if len(toolNames) > 0 {
		for _, name := range toolNames {
			delete(as.ActiveTools, name)
		}
	}

	if len(envNames) > 0 {
		for _, name := range envNames {
			delete(as.ActiveEnv, name)
		}
	}

	var newPaths []string
	for _, path := range as.ActivePaths {
		if arrutil.StringsContains(paths, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	as.ActivePaths = newPaths
}

// AddActivePath 添加激活路径, 会先检测是否已存在
func (as *ActivityState) AddActivePath(path string) {
	// 检查路径是否已存在
	for _, p := range as.ActivePaths {
		if p == path {
			return
		}
	}
	as.ActivePaths = append(as.ActivePaths, path)
}
