package kscript

import (
	"errors"

	"github.com/gookit/goutil/fsutil"
)

// AllowTypes shell wrapper for run script
var AllowTypes = []string{"sh", "zsh", "bash"}

// AllowExt list
var AllowExt = []string{".sh", ".zsh", ".bash", ".php", ".go", ".gop", ".kts", ".java", ".gry", ".groovy"}

// ExtToBinMap data
//
// eg:
//
//	'#!/usr/bin/env bash'
//	'#!/usr/bin/env -S go run'
var ExtToBinMap = map[string]string{
	".sh":     "sh",
	".zsh":    "zsh",
	".bash":   "bash",
	".php":    "php",
	".gry":    "groovy",
	".groovy": "groovy",
	".go":     "go run",
}

// NewRunner instance
func NewRunner(fns ...func(sr *Runner)) *Runner {
	sr := &Runner{
		ParseEnv:     true,
		AllowedExt:   AllowExt,
		ExtToBinMap:  ExtToBinMap,
		PathResolver: fsutil.ResolvePath,
		scriptFiles:  map[string]string{},
	}

	for _, fn := range fns {
		fn(sr)
	}
	return sr
}

// ErrNotFound error
var ErrNotFound = errors.New("script not found")

// IsNoNotFound error
func IsNoNotFound(err error) bool {
	return err != nil && err != ErrNotFound
}
