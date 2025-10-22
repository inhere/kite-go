package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gookit/goutil/arrutil"
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
	if runtime.GOOS == "windows" {
		return xenvHookShell == "bash" || strings.Contains(os.Getenv("SHELL"), "bash")
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

var winDiskPrefix = regexp.MustCompile(`^[a-zA-Z]:`)

// NormalizePath normalizes a path by expanding home directory and cleaning it
func NormalizePath(path string) string {
	fmtPath := filepath.Clean(fsutil.ExpandPath(path))

	if IsHookWinBash() {
		// Windows Git-Bash: 需要转换为 Unix 路径，同时需要处理盘符 eg: D:/ 转换为 /d/
		fmtPath = winDiskPrefix.ReplaceAllStringFunc(fsutil.UnixPath(fmtPath), func(sub string) string {
			return "/" + strings.ToLower(string(sub[0]))
		})
	}
	return fmtPath
}

// OutputScript outputs shell scripts to stdout
func OutputScript(script string) {
	if script != "" {
		fmt.Printf("%s\n%s\n", ScriptMark, script)
	}
}

// DiffRemovePaths diffs and removes paths from the PATH
func DiffRemovePaths(osPaths, rmPaths []string) (fmtRmPaths, newPaths, notFounds []string) {
	// format input paths
	for _, p := range rmPaths {
		fmtRmPaths = append(fmtRmPaths, NormalizePath(p))
	}

	var founds map[string]bool
	// find and remove from session PATH
	for _, p := range osPaths {
		if arrutil.StringsContains(fmtRmPaths, p) {
			founds[p] = true
		} else {
			newPaths = append(newPaths, p)
		}
	}

	// check found paths
	for _, p := range fmtRmPaths {
		if !founds[p] {
			notFounds = append(notFounds, p)
		}
	}
	return
}
