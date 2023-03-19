package kiteext

import (
	"strings"

	"github.com/gookit/goutil/maputil"
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
