package models

import (
	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite-go/pkg/util"
)

type ActivateSDKsParams struct {
	OpFlag OpFlag // 是否需要保存 direnv, global 文件
	// rem old active paths
	RemPaths []string
	// for  activate
	AddPaths []string
	AddEnvs  map[string]string
	AddSdks map[string]string
}

// NewActivateSDKsParams creates a new ActivateSDKsParams instance
func NewActivateSDKsParams() *ActivateSDKsParams {
	return &ActivateSDKsParams{
		AddEnvs: make(map[string]string),
		AddSdks: make(map[string]string),
		// AddTools: make(map[string]string),
	}
}

// IsGlobal 检测是否需要保存 global 文件
func (p *ActivateSDKsParams) IsGlobal() bool {
	return p.OpFlag == OpFlagGlobal
}

// IsDirenv 检测是否需要保存 direnv 文件
func (p *ActivateSDKsParams) IsDirenv() bool {
	return p.OpFlag == OpFlagDirenv
}

// AddSdk 添加激活工具链
func (p *ActivateSDKsParams) AddSdk(name, version string) {
	p.AddSdks[name] = version
}

// AddPath 添加激活路径, 会先检测是否已存在
func (p *ActivateSDKsParams) AddPath(path string) {
	for _, p := range p.AddPaths {
		if p == path {
			return
		}
	}
	p.AddPaths = append(p.AddPaths, path)
}

func (p *ActivateSDKsParams) AddSetEnvs(envs map[string]string) {
	p.AddEnvs = maputil.AppendSMap(p.AddEnvs, envs)
}

// AddRemPath 删除激活路径
func (p *ActivateSDKsParams) AddRemPath(path string) {
	for _, p := range p.RemPaths {
		if p == path {
			return
		}
	}
	p.RemPaths = append(p.RemPaths, path)
}

// GenInitScriptParams 生成初始化脚本参数
type GenInitScriptParams struct {
	// OpFlag OpFlag // 是否需要保存 direnv, global 文件
	// for  activate
	Paths []string
	Envs  map[string]string
	// add shell aliases
	ShellAliases map[string]string
	// ShellHooksDir shell hooks directory path
	ShellHooksDir string
}

// AddPath 添加环境PATH
func (p *GenInitScriptParams) AddPath(path string) {
	p.Paths = append(p.Paths, util.NormalizePath(path))
}

// AddPaths 添加环境PATH
func (p *GenInitScriptParams) AddPaths(paths []string) {
	for _, path := range paths {
		p.AddPath(path)
	}
}
