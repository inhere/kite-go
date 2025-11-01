package models

import (
	"os"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/fsutil"
)

type OpFlag uint8

const (
	OpFlagSession OpFlag = iota
	OpFlagDirenv
	OpFlagGlobal
)

func (of OpFlag) String() string {
	switch of {
	case OpFlagSession:
		return "session"
	case OpFlagDirenv:
		return "direnv"
	case OpFlagGlobal:
		return "global"
	}
	return "unknown"
}

// GetOpFlag 根据参数获取操作标识
func GetOpFlag(saveDirenv, global bool) OpFlag {
	if global {
		return OpFlagGlobal
	}
	if saveDirenv {
		return OpFlagDirenv
	}
	return OpFlagSession
}

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
		sessionID = time.Now().Format("20060102_150405")
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
	//
	// value 工具版本:
	//  - value 可以是 "*"
	//  - value 可以是 "*,required"
	//  - value 可以是 ">1.2"
	//  - value 可以是 ">1.2,required"
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
		Tools: make(map[string]string),
		Paths: []string{},
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}
}

// SessionID 从 as.File 获取当前会话ID. NOTE: 必须在 session 下使用
func (as *ActivityState) SessionID() string {
	return fsutil.NameNoExt(as.File)
}

// AddSDKs 新增激活工具
func (as *ActivityState) AddSDKs(sdks map[string]string) *ActivityState {
	for name, version := range sdks {
		as.SDKs[name] = version
	}
	return as
}

// AddEnvs 新增激活环境变量
func (as *ActivityState) AddEnvs(envs map[string]string) *ActivityState {
	for name, value := range envs {
		as.Envs[name] = value
	}
	return as
}

// AddTools 新增激活工具
func (as *ActivityState) AddTools(tools map[string]string) *ActivityState {
	for name, version := range tools {
		as.Tools[name] = version
	}
	return as
}

// AddPaths 新增激活路径
func (as *ActivityState) AddPaths(paths []string) *ActivityState {
	for _, path := range paths {
		as.AddPath(path)
	}
	return as
}

// Merge other to current. 合并两个状态数据
func (as *ActivityState) Merge(other *ActivityState) {
	if other == nil || other.IsEmpty() {
		return
	}
	as.AddSDKs(other.SDKs).AddEnvs(other.Envs).AddTools(other.Tools).AddPaths(other.Paths)
}

// DelSDKsEnvsPaths 删除激活工具和相关的 ENV, PATH
func (as *ActivityState) DelSDKsEnvsPaths(sdkNames, envNames, paths []string) {
	if len(sdkNames) > 0 {
		for _, name := range sdkNames {
			delete(as.SDKs, name)
		}
	}

	if len(envNames) > 0 {
		for _, name := range envNames {
			delete(as.Envs, name)
		}
	}

	if len(paths) > 0 {
		as.DelPaths(paths)
	}
}

// DelSDKs 删除多个SDK工具
func (as *ActivityState) DelSDKs(names []string) {
	for _, name := range names {
		delete(as.SDKs, name)
	}
}

// RemoveSDK 删除激活的SDK
func (as *ActivityState) RemoveSDK(name string) bool {
	_, exists := as.SDKs[name]
	if exists {
		delete(as.SDKs, name)
	}
	return exists
}

// DelTool 删除激活的工具
func (as *ActivityState) DelTool(name string) bool {
	_, exists := as.Tools[name]
	if exists {
		delete(as.Tools, name)
	}
	return exists
}

// DelThenAddPaths 先删除然后新增激活路径
func (as *ActivityState) DelThenAddPaths(rmPaths, addPaths []string) *ActivityState {
	return as.DelPaths(rmPaths).AddPaths(addPaths)
}

// DelPaths 删除激活路径
func (as *ActivityState) DelPaths(paths []string) *ActivityState {
	if len(paths) == 0 {
		return as
	}

	var newPaths []string
	for _, path := range as.Paths {
		if arrutil.StringsContains(paths, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	as.Paths = newPaths
	return as
}

// AddPath 添加激活路径, 会先检测是否已存在
func (as *ActivityState) AddPath(path string) {
	// 检查路径是否已存在
	for _, p := range as.Paths {
		if p == path {
			return
		}
	}
	as.Paths = append(as.Paths, path)
}

// ExistsPath 检查路径是否已存在
func (as *ActivityState) ExistsPath(val string) bool {
	return arrutil.StringsContains(as.Paths, val)
}

// IsEmpty 检查状态数据是否为空
func (as *ActivityState) IsEmpty() bool {
	return len(as.SDKs) == 0 &&
		len(as.Envs) == 0 &&
		len(as.Paths) == 0 &&
		len(as.Tools) == 0
}
