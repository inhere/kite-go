package util

import (
	"os"
	"runtime"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/sysutil"
)

var xenvHookShell = envutil.Getenv("XENV_HOOK_SHELL")

// HookShell returns the hook shell name. 不为空表明在shell hook环境中
func HookShell() string { return xenvHookShell }

// InHookShell returns true if the current shell is in the hook shell
func InHookShell() bool { return xenvHookShell != "" }

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

// ClinkIsInstalled checks if Clink is installed on Windows
func ClinkIsInstalled() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return sysutil.HasExecutable("clink.exe")
}

// ToExecutableName converts a tool name to the executable name depending on OS
func ToExecutableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
