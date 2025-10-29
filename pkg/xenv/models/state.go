package models

import (
	"os"
	"time"

	"github.com/gookit/goutil/arrutil"
)

const (
	// GlobalStateFile global state file path
	GlobalStateFile = "~/.config/xenv/global.toml"
	// LocalStateFile local state file path
	LocalStateFile = ".xenv.toml"
	// SessIdEnvName 当前会话ID环境变量名称
	SessIdEnvName = "XENV_SESSION_ID"
	// SessionStateDir 当前SHELL会话状态文件目录 eg: ~/.xenv/session/<session_id>.json
	SessionStateDir = "~/.xenv/session"
)

var sessionID = os.Getenv(SessIdEnvName)

// SessionID 获取当前会话ID
func SessionID() string {
	if sessionID == "" {
		sessionID = time.Now().Format("20060102150405")
	}
	return sessionID
}

// SessionStateFile 生成当前会话状态文件
func SessionStateFile() string {
	return SessionStateDir + "/" + SessionID() + ".json"
}

// ActivityState 代表用户当前激活的工具链和环境状态.
//  - 全局的会保存到 ~/.config/xenv/global.toml
//  - 目录级的会保存到 {pwd|parent}/.xenv.toml
//  - 会话的会保存到 ~/.xenv/session/<session_id>.json
type ActivityState struct {
	// 激活的 ENV 路径列表
	Paths []string `json:"paths" toml:"paths"`
	// 激活的工具链映射 key为工具名，value为版本 Sdks
	SDKs map[string]string `json:"sdks" toml:"sdks"`
	// 激活的环境变量
	Envs map[string]string `json:"envs" toml:"envs"`
	// Tools 需要的工具列表
	Tools map[string]string `json:"tools" toml:"tools"`
	// state file path OR ID string.
	File string `json:"-" toml:"-"`
	// 是否有更新 - 内部使用，用于标识状态数据是否需要更新
	HasUpdate bool `json:"-" toml:"-"`
	// Shell 当前使用的shell - 仅在 session 下有效
	Shell string `json:"shell" toml:"-"`

	// EnableGlobal 是否启用全局环境配置
	// EnableGlobal bool `json:"enable_global"`
	// 创建时间
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

// NewActivityState creates a new ActivityState
func NewActivityState(filePath string) *ActivityState {
	return &ActivityState{
		File: filePath,
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

// ExistsPath 检查路径是否已存在
func (as *ActivityState) ExistsPath(val string) bool {
	return arrutil.StringsContains(as.Paths, val)
}
