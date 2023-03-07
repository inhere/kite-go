package kiteext

import (
	"runtime"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// PathMap struct
type PathMap struct {
	maputil.Aliases `json:"pathmap"`
}

// Resolve path
func (p *PathMap) Resolve(path string) string {
	if idx := strings.IndexByte(path, '/'); idx > 0 {
		prefix, other := path[:idx], path[idx:]
		return p.ResolveAlias(prefix) + other
	}

	return p.ResolveAlias(path)
}

// Data map get
func (p *PathMap) Data() map[string]string {
	return p.Aliases
}

// VarMap struct
type VarMap struct {
	maputil.Aliases
	Prefix string

	replacer *strings.Replacer
}

// NewVarMap instance
func NewVarMap(smp map[string]string) *VarMap {
	vm := &VarMap{
		Prefix: "$",
	}

	vm.LoadMap(map[string]string{
		"os":   runtime.GOOS,
		"user": sysutil.CurrentUser().Name,
	})

	vm.LoadMap(smp)
	return vm
}

// LoadMap path
func (m *VarMap) LoadMap(smp map[string]string) {
	if smp != nil {
		for k, v := range smp {
			m.Aliases[m.Prefix+k] = v
		}
		m.replacer = strutil.NewReplacer(m.Aliases)
	}
}

// Replace vars in string.
func (m *VarMap) Replace(s string) string {
	if strings.Contains(s, m.Prefix) {
		return m.replacer.Replace(s)
	}
	return s
}

// Resolve path
func (m *VarMap) Resolve(path string) string {
	if idx := strings.IndexByte(path, '/'); idx > 0 {
		prefix, other := path[:idx], path[idx:]
		return m.ResolveAlias(prefix) + other
	}

	return m.ResolveAlias(path)
}

// Data map get
func (m *VarMap) Data() map[string]string {
	return m.Aliases
}
