package shell

import (
	"runtime"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/sysutil"
)

// This file contains shell integration utilities

var xenvHookShell = envutil.Getenv("XENV_HOOK_SHELL")

// HookShell returns the hook shell name
func HookShell() string { return xenvHookShell }

// InHookShell returns true if the current shell is in the hook shell
func InHookShell() bool { return xenvHookShell != "" }

// IsValidShellType checks if a shell type is valid
func IsValidShellType(shellType string) bool {
	switch shellType {
	case string(Bash), string(Zsh), string(Pwsh), string(Cmd):
		return true
	default:
		return false
	}
}

// TypeFromString returns the shell type from a string
func TypeFromString(shellType string) ShellType {
	shellType = strings.ToLower(shellType)
	switch shellType {
	case "bash":
		return Bash
	case "zsh":
		return Zsh
	case "pwsh", "powershell":
		return Pwsh
	default:
		return Cmd
	}
}

// ClinkIsInstalled checks if Clink is installed on Windows
func ClinkIsInstalled() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return sysutil.HasExecutable("clink.exe")
}
