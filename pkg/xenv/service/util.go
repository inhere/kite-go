package service

import (
	"github.com/inhere/kite-go/pkg/util"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
)

// getShellGenerator 获取当前shell的脚本生成器. 注意：不在shell hook环境，会返回nil
func getShellGenerator(_ *models.Configuration) (*shell.XenvScriptGenerator, error) {
	// hookShell 不为空表明在shell hook环境中
	hookShell := util.HookShell()
	if hookShell == "" {
		return nil, nil
	}

	shellType, err := shell.TypeFromString(hookShell)
	if err != nil {
		return nil, err
	}

	return shell.NewScriptGenerator(shellType), nil
}

