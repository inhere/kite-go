package util

import (
	"runtime"

	"github.com/gookit/goutil/sysutil"
)


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
