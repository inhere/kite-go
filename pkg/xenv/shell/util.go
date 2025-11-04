package shell

import (
	"fmt"

	"github.com/gookit/goutil/arrutil"
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

// This file contains shell integration utilities

// OutputScript outputs shell scripts to stdout
func OutputScript(script string) {
	if script != "" {
		// fix: Windows Pwsh, CMD 换行符是 \r\n
		if xenvcom.IsHookPwshOrCmd() {
			fmt.Printf("%s\r\n%s\r\n", ScriptMark, script)
			return
		}
		fmt.Printf("%s\n%s\n", ScriptMark, script)
	}
}

// DiffRemovePaths diffs and removes paths from the PATH
func DiffRemovePaths(osPaths, rmPaths []string) (fmtRmPaths, newPaths, notFounds []string) {
	if len(rmPaths) == 0 {
		return nil, osPaths, nil
	}

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
