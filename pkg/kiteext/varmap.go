package kiteext

import (
	"runtime"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// VarMap struct
type VarMap struct {
	Prefix   string
	aliases maputil.Aliases
	// replacer build from aliases data
	replacer *strings.Replacer
	// raw vars map, key not append Prefix
	rawVars map[string]string
}

// NewVarMap instance
func NewVarMap(smp map[string]string) *VarMap {
	vm := &VarMap{
		Prefix:  "$",
		aliases: make(maputil.Aliases),
	}

	vm.LoadMap(smp)
	vm.LoadMap(map[string]string{
		"os":   runtime.GOOS,
		"user": sysutil.CurrentUser().Name,
	})
	return vm
}

// LoadMap path
func (m *VarMap) LoadMap(smp map[string]string) {
	if len(smp) == 0 {
		return
	}

	// update alias map
	for k, v := range smp {
		m.aliases[m.Prefix+k] = v
	}

	m.rawVars = maputil.MergeStringMap(smp, m.rawVars, false)
}

// Add new var value
func (m *VarMap) Add(name, value string) {
	m.rawVars[name] = value
	m.aliases[m.Prefix+name] = value
}

func (m *VarMap) ensureReplacer() {
	if m.replacer == nil {
		m.replacer = strutil.NewReplacer(m.aliases)
	}
}

// Replace vars in string.
func (m *VarMap) Replace(s string) string {
	if strings.Contains(s, m.Prefix) {
		m.ensureReplacer()
		return m.replacer.Replace(s)
	}
	return s
}

// Resolve a string value.
func (m *VarMap) Resolve(path string) string {
	// eg: $home/some.txt
	if idx := strings.IndexByte(path, '/'); idx > 0 {
		first, other := path[:idx], path[idx:]
		return m.aliases.ResolveAlias(first) + other
	}

	return m.aliases.ResolveAlias(path)
}

// AliasMap data get. key is prefix + alias
func (m *VarMap) AliasMap() map[string]string { return m.aliases }

// Data raw map get. key is without prefix
func (m *VarMap) Data() map[string]string { return m.rawVars }
