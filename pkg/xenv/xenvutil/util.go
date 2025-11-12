package xenvutil

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/util"
)

// ListVersionDirs 列出SDK的已安装版本目录
//
// return key: version, value: dir path
func ListVersionDirs(installDir string) (map[string]string, error) {
	// 获取SDK基础目录
	baseDir := filepath.Dir(util.NormalizePath(installDir))
	if !fsutil.IsDir(baseDir) {
		return nil, fmt.Errorf("SDK directory does not exist: %s", baseDir)
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

// ParseGoVersion 实现简单的从 go.mod 文件解析go版本
//  - 按行读取文件 找到 go {version} 所在行即停止，最多读取 10 行，找不到就返回
//
// eg: goVer, err := ParseGoVersion("go.mod")
func ParseGoVersion(modFile string) (string, error) {
	file, err := os.Open(modFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	goPrefix := "go "

	for scanner.Scan() && lineCount < 10 {
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		// eg: go 1.19
		if strings.HasPrefix(line, goPrefix) {
			version := strings.TrimSpace(line[len(goPrefix):])
			if version != "" {
				return version, nil
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("go version not found in first 10 lines")
}

// ParseToolVersions 实现从 .tool-versions 文件解析工具版本
func ParseToolVersions(toolFile string) (map[string]string, error) {
	contents, err := os.ReadFile(toolFile)
	if err != nil {
		return nil, err
	}

	// 解析 .tool-versions 文件内容
	lines := strings.Split(string(contents), "\n")
	versions := make(map[string]string)

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		toolName := strings.TrimSpace(parts[0])
		toolVersion := strings.TrimSpace(parts[1])
		if toolName == "" || toolVersion == "" {
			continue
		}
		versions[toolName] = toolVersion
	}

	return versions, nil
}

// ParseNvmrcFile 实现从 .nvmrc 文件解析工具版本
//
//  - 纯文本，只包含一个 Node.js 版本号
//  - 内容可能为：18.17.0, v14.18.1, lts/*, node(最新稳定版: node)
func ParseNvmrcFile(nvmrcFile string) (string, error) {
	contents, err := os.ReadFile(nvmrcFile)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(contents))
	if version == "" {
		return "", fmt.Errorf("invalid .nvmrc file: %s", nvmrcFile)
	}

	version = strings.Trim(version, "v/*")
	// TODO 暂时不支持检查 node 最新稳定版
	if version == "node" {
		version = "latest"
	}
	return version, nil
}

// ParsePythonVersion 实现从 .python-version 文件解析 Python 版本
//  - 纯文本，只包含一个 Python
//  - 内容可能为：3.11.4,
func ParsePythonVersion(versionFile string) (string, error) {
	contents, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(contents))
	if version == "" {
		return "", fmt.Errorf("invalid .python-version file: %s", versionFile)
	}

	version = strings.Trim(version, "v*")
	return version, nil
}
