package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
)

// getShellGenerator 获取当前shell的脚本生成器. 注意：不在shell hook环境，会返回nil
func getShellGenerator(_ *models.Configuration) (*shell.XenvScriptGenerator, error) {
	// hookShell 不为空表明在shell hook环境中
	hookShell := xenvcom.HookShell()
	if hookShell == "" {
		return nil, nil
	}

	shellType, err := shell.TypeFromString(hookShell)
	if err != nil {
		return nil, err
	}

	return shell.NewScriptGenerator(shellType), nil
}

// parseGoVersion 实现简单的从 go.mod 文件解析go版本
//  - 按行读取文件 找到 go {version} 所在行即停止，最多读取 10 行，找不到就返回
//
// eg: goVer, err := parseGoVersion("go.mod")
func parseGoVersion(modFile string) (string, error) {
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

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("go version not found in first 10 lines")
}

// parseToolVersions 实现从 .tool-versions 文件解析工具版本
func parseToolVersions(toolFile string) (map[string]string, error) {
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

// parseNvmrcFile 实现从 .nvmrc 文件解析工具版本
//
//  - 纯文本，只包含一个 Node.js 版本号
//  - 内容可能为：18.17.0, v14.18.1, lts/*, node(最新稳定版: node)
func parseNvmrcFile(nvmrcFile string) (string, error) {
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

// parsePythonVersion 实现从 .python-version 文件解析 Python 版本
//  - 纯文本，只包含一个 Python
//  - 内容可能为：3.11.4,
func parsePythonVersion(versionFile string) (string, error) {
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
