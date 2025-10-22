package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gookit/goutil/fsutil"
)

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

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return fsutil.EnsureDir(NormalizePath(path))
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// DirExists checks if a directory exists
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

// CreateSymlink creates a symbolic link
func CreateSymlink(target, linkPath string) error {
	// Check if the link already exists
	exists, err := FileExists(linkPath)
	if err != nil {
		return err
	}

	if exists {
		// Remove existing link/file
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing file: %w", err)
		}
	}

	return os.Symlink(target, linkPath)
}

// FindExecutable looks for an executable file in the given directories
func FindExecutable(name string, paths []string) (string, error) {
	for _, path := range paths {
		fullPath := filepath.Join(path, name)
		if exists, err := FileExists(fullPath); err == nil && exists {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("executable %s not found in provided paths", name)
}
