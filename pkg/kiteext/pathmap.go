package kiteext

import (
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

// PathMap struct
type PathMap struct {
	// internal, split from Prefixes
	prefixes []string
	// path alias map: {alias: path}
	maputil.Aliases `json:"path_map"`
	// if true, must start with prefix on resolve path.
	Strict bool `json:"strict"`
	// allowed prefix mark for alias, use on resolve. eg: "$,#,@"
	Prefixes string `json:"prefixes"`
	// FallbackFn if alias not exists, will call this func to get path.
	FallbackFn func(path string) string
}

// NewPathMap instance
func NewPathMap(opFns ...func(pm *PathMap)) *PathMap {
	pm := &PathMap{
		Aliases:  make(maputil.Aliases),
		Prefixes: "$,#,@",
	}

	for _, fn := range opFns {
		fn(pm)
	}
	return pm
}

// Resolve path. eg: $home/some.txt => /home/USER/some.txt
func (p *PathMap) Resolve(path string) string {
	if len(path) == 0 {
		return path
	}

	rawPath := path
	if p.prefixes == nil {
		p.prefixes = strutil.Split(p.Prefixes, ",")
	}

	// check prefix mark
	var hasMark bool
	if len(p.prefixes) > 0 {
		hasMark = strutil.HasOnePrefix(path, p.prefixes)
		if p.Strict && !hasMark {
			return path
		}

		if hasMark {
			path = path[1:]
		}
	}

	if !hasMark && fsutil.PathExists(path) {
		return path
	}

	alias := path
	// check sub path
	var other string
	if idx := strings.IndexByte(path, '/'); idx > 0 {
		alias, other = path[:idx], path[idx:]
	}

	// resolve alias
	if p.HasAlias(alias) {
		return p.ResolveAlias(alias) + other
	}
	if p.FallbackFn != nil {
		return p.FallbackFn(rawPath)
	}
	return path
}

// Data get all alias map data
func (p *PathMap) Data() map[string]string {
	return p.Aliases
}
