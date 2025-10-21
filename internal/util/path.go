package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ExpandHome expands the tilde (~) to the user's home directory
func ExpandHome(path string) string {
	// Check if the path starts with ~
	if path == "~" {
		return HomeDir()
	}

	if strings.HasPrefix(path, "~/") {
		homeDir := HomeDir()
		return filepath.Join(homeDir, path[2:])
	}

	return path
}

// HomeDir returns the user's home directory
func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// NormalizePath normalizes a path by expanding home directory and cleaning it
func NormalizePath(path string) string {
	expanded := ExpandHome(path)
	return filepath.Clean(expanded)
}

// ToExecutableName converts a tool name to the executable name depending on OS
func ToExecutableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
