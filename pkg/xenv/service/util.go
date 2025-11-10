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
