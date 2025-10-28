package shell

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/envutil"
	"github.com/inhere/kite-go/pkg/util"
)

// This file contains shell integration utilities

var xenvHookShell = envutil.Getenv("XENV_HOOK_SHELL")

// HookShell returns the hook shell name. 不为空表明在shell hook环境中
func HookShell() string { return xenvHookShell }

// InHookShell returns true if the current shell is in the hook shell
func InHookShell() bool { return xenvHookShell != "" }

// IsHookBash checks if the current hook shell is Linux/Windows Bash(eg: git-bash)
func IsHookBash() bool {
	if runtime.GOOS == "windows" {
		return xenvHookShell == "bash" || strings.Contains(os.Getenv("SHELL"), "bash")
	}
	return false
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
		fmtRmPaths = append(fmtRmPaths, util.NormalizePath(p))
	}

	founds := make(map[string]bool)
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
