package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IsWindows returns true if the current platform is Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

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
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to HOME or USERPROFILE environment variable
		if runtime.GOOS == "windows" {
			home = os.Getenv("USERPROFILE")
		} else {
			home = os.Getenv("HOME")
		}
	}
	return home
}

// NormalizePath normalizes a path by expanding home directory and cleaning it
func NormalizePath(path string) string {
	expanded := ExpandHome(path)
	return filepath.Clean(expanded)
}

// AddToPath adds a directory to the PATH environment variable
func AddToPath(dir string) error {
	// In a real implementation, this would modify the PATH for the current session
	// For now, this function would handle the logic of adding a path
	// This is typically done in shell hooks, not directly in Go
	return nil
}

// GetPathSeparator returns the appropriate path separator for the current OS
func GetPathSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	}
	return ":"
}

// JoinPathList joins multiple path entries into a single PATH string
func JoinPathList(paths []string) string {
	return strings.Join(paths, GetPathSeparator())
}

// SplitPathList splits a PATH string into individual paths
func SplitPathList(pathList string) []string {
	return strings.Split(pathList, GetPathSeparator())
}

// ConvertToExecutableName converts a tool name to the executable name depending on OS
func ConvertToExecutableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}