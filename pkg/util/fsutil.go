package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

var winDiskPrefix = regexp.MustCompile(`^[a-zA-Z]:`)

// NormalizePath normalizes a path by expanding home directory and cleaning it
func NormalizePath(path string) string {
	fmtPath := filepath.Clean(fsutil.ExpandPath(path))
	if xenvcom.IsHookBash() {
		fmtPath = fsutil.UnixPath(fmtPath)
	}
	return fmtPath
}

// FmtEnvPath formats an environment path for use in the current shell
func FmtEnvPath(envPath string) string {
	if xenvcom.IsHookBash() {
		envPath = fsutil.UnixPath(envPath)
		// Windows Git-Bash: 需要转换为 Unix 路径，同时需要处理盘符 eg: D:/ 转换为 /d/
		envPath = winDiskPrefix.ReplaceAllStringFunc(envPath, func(sub string) string {
			return "/" + strings.ToLower(string(sub[0]))
		})
	}
	return envPath
}

// SplitPath splits a PATH string into individual paths
func SplitPath(envPath string) []string {
	// NOTE: 分割是需要使用 os.PathListSeparator
	return strings.Split(envPath, string(os.PathListSeparator))
}

// JoinPaths joins multiple path entries into a single PATH string
func JoinPaths(paths []string) string {
	return strings.Join(paths, xenvcom.PathSep())
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
		if err1 := os.Remove(linkPath); err1 != nil {
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
