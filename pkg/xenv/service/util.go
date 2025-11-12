package service

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/pkg/xenv/models"
	"github.com/inhere/kite-go/pkg/xenv/shell"
	"github.com/inhere/kite-go/pkg/xenv/xenvcom"
	"github.com/inhere/kite-go/pkg/xenv/xenvutil"
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

func sdkVersionsFromSpecifiedFiles(specMap map[string]*models.VersionSpec) {
	// 支持识别常用的工具配置 eg: go.mod, .tool-versions, .nvmrc, .python-version
	toolsCfgFiles := []string{"go.mod", ".tool-versions", ".nvmrc", ".python-version"}
	for _, filename := range toolsCfgFiles {
		if !fsutil.IsFile(filename) {
			continue
		}

		switch filename {
		case ".tool-versions":
			// 识别 .tool-versions 文件
			verMap, err := xenvutil.ParseToolVersions(filename)
			if err != nil {
				ccolor.Warnf("Failed to parse .tool-versions file: %v\n", err)
				continue
			}

			ccolor.Infof("Detect tool versions from .tool-versions: %v\n", verMap)
			for name, ver := range verMap {
				specMap[name] = &models.VersionSpec{
					Name:    name,
					Version: ver,
				}
			}
		case "go.mod":
			goVer, err := xenvutil.ParseGoVersion(filename)
			if err != nil {
				ccolor.Warnf("Failed to parse go.mod file: %v\n", err)
				continue
			}
			ccolor.Infof("Detect go version from go.mod: %s\n", goVer)
			specMap["go"] = &models.VersionSpec{
				Name:    "go",
				Version: goVer,
			}
		case ".nvmrc":
			// 识别 .nvmrc 文件
			nodeVer, err := xenvutil.ParseNvmrcFile(filename)
			if err != nil {
				ccolor.Warnf("Failed to parse .nvmrc file: %v\n", err)
				continue
			}
			ccolor.Infof("Detect node version from .nvmrc: %s\n", nodeVer)
			specMap["node"] = &models.VersionSpec{
				Name:    "node",
				Version: nodeVer,
			}
		case ".python-version":
			// 识别 .python-version 文件
			pyVer, err := xenvutil.ParsePythonVersion(filename)
			if err != nil {
				ccolor.Warnf("Failed to parse .python-version file: %v\n", err)
				continue
			}
			ccolor.Infof("Detect python version from .python-version: %s\n", pyVer)
			specMap["python"] = &models.VersionSpec{
				Name:    "python",
				Version: pyVer,
			}
		}
	}
}

