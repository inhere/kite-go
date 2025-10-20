package shell

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gookit/goutil/envutil"
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
