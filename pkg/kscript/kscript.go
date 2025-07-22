package kscript

import (
	"errors"

	"github.com/gookit/goutil/fsutil"
)

var (
	// AllowTypes shell wrapper for run a script. eg: bash, sh, zsh, cmd, pwsh
	AllowTypes = []string{"sh", "zsh", "bash", "cmd", "pwsh"}

	// AllowExt list. allowed script file ext.
	AllowExt = []string{".sh", ".zsh", ".bash", ".php", ".go", ".gop", ".kts", ".java", ".gry", ".groovy", ".py"}
)

var (
	// DefaultTaskFiles 默认自动查找的task文件名称 eg "kite.task[s].yml", "kite.script[s].yml"
	DefaultTaskFiles = []string{".kite.task", ".kite.tasks", ".kite.script", ".kite.scripts"}
	// DefaultDefineExts 默认允许的 scriptApp, scriptTask 定义文件后缀
	DefaultDefineExts = []string{".yml", ".yaml", ".toml", ".json"}
)

// ExtToBinMap data
//
// eg:
//
//	'#!/usr/bin/env bash'
//	'#!/usr/bin/env -S go run'
var ExtToBinMap = map[string]string{
	".sh":   "sh",
	".zsh":  "zsh",
	".bash": "bash",
	".php":  "php",
	".py":   "python",
	// ".dart":   "dart",
	".gry":    "groovy",
	".groovy": "groovy",
	".go":     "go run",
}

// NewRunner instance
func NewRunner(fns ...func(kr *Runner)) *Runner {
	kr := &Runner{
		ParseEnv:     true,
		PathResolver: fsutil.ResolvePath,
		// script file
		AllowedExt:   AllowExt,
		ExtToBinMap:  ExtToBinMap,
		scriptFiles:  map[string]string{},
		// script app
		ScriptApps:    []string{"?$base/script-app"},
		ScriptAppExts: DefaultDefineExts,
		// script task
		AutoTaskFiles: DefaultTaskFiles,
		AutoTaskExts:  DefaultDefineExts,
		AutoMaxDepth:  6,
	}

	for _, fn := range fns {
		fn(kr)
	}
	return kr
}

// ErrNotFound error
var ErrNotFound = errors.New("script not found")

// IsNoNotFound error
func IsNoNotFound(err error) bool {
	return err != nil && !errors.Is(err, ErrNotFound)
}
