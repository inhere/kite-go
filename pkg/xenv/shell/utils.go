package shell

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
)

// This file contains shell integration utilities

var xenvHookShell = envutil.Getenv("XENV_HOOK_SHELL")

// HookShell returns the hook shell name. 不为空表明在shell hook环境中
func HookShell() string { return xenvHookShell }

// InHookShell returns true if the current shell is in the hook shell
func InHookShell() bool { return xenvHookShell != "" }

// IsHookWinBash checks if the current hook shell is Windows Bash(eg: git-bash)
func IsHookWinBash() bool {
	return runtime.GOOS == "windows" && xenvHookShell == "bash"
}

// IsValidShellType checks if a shell type is valid
func IsValidShellType(shellType string) bool {
	shellType = strings.ToLower(shellType)
	switch shellType {
	case string(Bash), string(Zsh), string(Pwsh), string(Cmd):
		return true
	default:
		return false
	}
}

// TypeFromString returns the shell type from a string
func TypeFromString(shellType string) (ShellType, error) {
	shellType = strings.ToLower(shellType)
	switch shellType {
	case "bash":
		return Bash, nil
	case "zsh":
		return Zsh, nil
	case "pwsh", "powershell":
		return Pwsh, nil
	case "cmd", "clink":
		return Cmd, nil
	default:
		return "", fmt.Errorf("unsupported shell type: %s (should: bash, zsh, pwsh or cmd)", shellType)
	}
}

// ClinkIsInstalled checks if Clink is installed on Windows
func ClinkIsInstalled() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return sysutil.HasExecutable("clink.exe")
}

// PathSeparator returns the appropriate path separator for the current OS
func PathSeparator() string {
	if runtime.GOOS == "windows" {
		if xenvHookShell == "bash" {
			return ":"
		}
		return ";"
	}
	return ":"
}

// SplitPath splits a PATH string into individual paths
func SplitPath(envPath string) []string {
	return strings.Split(envPath, PathSeparator())
}

// JoinPaths joins multiple path entries into a single PATH string
func JoinPaths(paths []string) string {
	return strings.Join(paths, PathSeparator())
}

// NormalizePath normalizes a path by expanding home directory and cleaning it
func NormalizePath(path string) string {
	return filepath.Clean(fsutil.ExpandPath(path))
}
