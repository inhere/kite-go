package models

import (
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

	//
	// 下面的字段 仅在 session 下有效
	//

	// Shell 当前使用的shell type
	Shell string `json:"shell,omitempty" toml:"-"`
	// 当前会话关联的所有目录状态数据. 用于跳转目录时，销毁之前的目录state
	//  - key: state file path, value: state data
	DirStates map[string]*ActivityState `json:"dir_states,omitempty" toml:"-"`

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

// IsSession 检查当前状态数据是否为会话状态
func (as *ActivityState) IsSession() bool {
	return as.Shell != ""
}

// SessionID 从 as.File 获取当前会话ID. NOTE: 必须在 session 下使用
func (as *ActivityState) SessionID() string {
	if as.Shell == "" {
		panic("state: session shell can not be empty")
	}
	return fsutil.NameNoExt(as.File)
}

// AddSDKs 新增激活工具
func (as *ActivityState) AddSDKs(sdks map[string]string) *ActivityState {
	for name, version := range sdks {
		as.SDKs[name] = version
	}
	as.HasUpdate = true
	return as
}

// AddEnvs 新增激活环境变量
func (as *ActivityState) AddEnvs(envs map[string]string) *ActivityState {
	for name, value := range envs {
		as.Envs[name] = value
	}
	as.HasUpdate = true
	return as
}

// AddTools 新增激活工具
func (as *ActivityState) AddTools(tools map[string]string) *ActivityState {
	for name, version := range tools {
		as.Tools[name] = version
	}
	as.HasUpdate = true
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
		as.DelSDKs(sdkNames)
	}
	if len(envNames) > 0 {
		as.DelEnvs(envNames)
	}
	if len(paths) > 0 {
		as.DelPaths(paths)
	}
}

// DelSDKs 删除多个SDK工具
func (as *ActivityState) DelSDKs(names []string) {
	for _, name := range names {
		as.HasUpdate = true
		delete(as.SDKs, name)
	}
}

// RemoveSDK 删除激活的SDK
func (as *ActivityState) RemoveSDK(name string) bool {
	_, exists := as.SDKs[name]
	if exists {
		as.HasUpdate = true
		delete(as.SDKs, name)
	}
	return exists
}

// DelEnvs 删除多个环境变量
func (as *ActivityState) DelEnvs(names []string) {
	for _, name := range names {
		as.HasUpdate = true
		delete(as.Envs, name)
	}
}

// DelTool 删除激活的工具
func (as *ActivityState) DelTool(name string) bool {
	_, exists := as.Tools[name]
	if exists {
		as.HasUpdate = true
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
		as.HasUpdate = true
		newPaths = append(newPaths, path)
	}
	as.Paths = newPaths
	return as
}

// DelPath 删除激活路径
func (as *ActivityState) DelPath(path string) {
	as.DelPaths([]string{path})
}

// AddPath 添加激活路径, 会先检测是否已存在
func (as *ActivityState) AddPath(path string) bool {
	// 检查路径是否已存在
	for _, p := range as.Paths {
		if p == path {
			return false
		}
	}
	as.HasUpdate = true
	as.Paths = append(as.Paths, path)
	return true
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
