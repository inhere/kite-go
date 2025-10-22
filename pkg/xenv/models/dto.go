package models

import "github.com/gookit/goutil/maputil"

type ActivateToolsParams struct {
	// rem old active paths
	RemPaths []string
	// for  activate
	AddPaths []string
	AddEnvs  map[string]string
	AddTools  map[string]string
}

// NewActivateToolsParams creates a new ActivateToolsParams instance
func  NewActivateToolsParams() *ActivateToolsParams {
	return &ActivateToolsParams{
		AddEnvs:  make(map[string]string),
		AddTools: make(map[string]string),
	}
}

// AddTool 添加激活工具
func (p *ActivateToolsParams) AddTool(name, version string) {
	p.AddTools[name] = version
}

// AddPath 添加激活路径, 会先检测是否已存在
func (p *ActivateToolsParams) AddPath(path string) {
	for _, p := range p.AddPaths {
		if p == path {
			return
		}
	}
	p.AddPaths = append(p.AddPaths, path)
}

func (p *ActivateToolsParams) AddSetEnvs(envs map[string]string) {
	p.AddEnvs = maputil.AppendSMap(p.AddEnvs, envs)
}

// AddRemPath 删除激活路径
func (p *ActivateToolsParams) AddRemPath(path string) {
	for _, p := range p.RemPaths {
		if p == path {
			return
		}
	}
	p.RemPaths = append(p.RemPaths, path)
}
