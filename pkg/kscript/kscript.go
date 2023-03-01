package kscript

import "github.com/gookit/goutil/sysutil"

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
		PathResolver: sysutil.ExpandPath,
		scriptFiles:  map[string]string{},
	}

	for _, fn := range fns {
		fn(sr)
	}
	return sr
}
