package xenvcom

import (
	"os"
	"runtime"
	"strings"

	"github.com/gookit/goutil/envutil"
)

var xenvHookShell = envutil.Getenv(HookShellEnvName)

// HookShell returns the hook shell name. 不为空表明在shell hook环境中
func HookShell() string { return xenvHookShell }

// InHookShell returns true if the current shell is in the hook shell
func InHookShell() bool { return xenvHookShell != "" }

// SetHookShell sets the hook shell name. NOTE: use for testing only
func SetHookShell(shell string) { xenvHookShell = shell }

// IsHookBash checks if the current hook shell is Windows Bash(eg: git-bash)
func IsHookBash() bool {
	if runtime.GOOS == "windows" {
		return xenvHookShell == "bash" || strings.Contains(os.Getenv("SHELL"), "bash")
	}
	return false
}

// IsHookPwshOrCmd checks if the current hook shell is Windows PowerShell or CMD
func IsHookPwshOrCmd() bool {
	if runtime.GOOS == "windows" {
		return xenvHookShell == "pwsh" || xenvHookShell == "cmd"
	}
	return false
}

// PathSep returns the appropriate path separator for the current OS
func PathSep() string {
	if runtime.GOOS == "windows" {
		if xenvHookShell == "bash" {
			return ":"
		}
		return ";"
	}
	return ":"
}
