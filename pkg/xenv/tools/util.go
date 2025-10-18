package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
)

// ListVersionDirs 列出SDK的已安装版本目录
//
// return key: version, value: dir path
func ListVersionDirs(installDir string) (map[string]string, error) {
	// 获取SDK基础目录
	baseDir := filepath.Dir(installDir)
	if !fsutil.IsDir(baseDir) {
		return nil, nil
	}
	// ccolor.Infof("DEBUG list installed version from %s\n", baseDir)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list tool directory: %w", err)
	}

	var sdkDirMap = make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name() // 从目录名中提取版本号
			if verStr := strutil.NumVersion(dirName); verStr != "" {
				sdkDirMap[verStr] = baseDir + "/" + entry.Name()
			}
		}
	}
	return sdkDirMap, nil
}
