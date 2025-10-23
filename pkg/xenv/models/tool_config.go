package models

import (
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

// ToolChain SDK 开发工具（如Go、Node.js等）配置，包含安装路径、别名等属性。
//  - SDKToolConfig, ToolConfig
//  - SDK信息配置，不含有特定的版本信息(支持多版本存在)
type ToolChain struct {
	Name  string `json:"name"`  // 工具名称，如 "go", "node", jdk
	Alias string `json:"alias"` // 工具别名列表，如 "golang" for go, jdk for java
	// 可选，下载URL模板 eg: "https://golang.org/dl/go{version}.{os}-{arch}.{download_ext}"
	InstallURL string `json:"install_url"`
	// 从远程下载不同OS平台的工具包的后缀格式
	// eg:
	//
	// 	download_ext:
	// 	  windows: zip
	// 	  linux: tar.gz
	// 	  darwin: tar.gz
	DownloadExt map[string]string `json:"download_ext"`
	// sdk tool 安装目录路径 默认 ~/.xenv/tools/{Name}/{version}
	//  - {version} 是动态的，根据版本号替换
	//  - 可以自定义 eg: ~/.xenv/tools/go/go{version}
	InstallDir string `json:"install_dir"`
	// 激活时设置的额外环境变量
	//  - value 可用使用一些内部变量，如 {version}，{install_dir}
	ActiveEnv map[string]string `json:"active_env"`
	// 该工具的 bin 文件目录名称，不设置就是 install_dir/bin 目录
	BinDir      string   `json:"bin_dir"`
	BinPaths    []string `json:"bin_paths"`    // 该工具提供的二进制文件路径列表
	// 安装完成后执行的shell hook脚本
	PostInstall []string `json:"post_install"`
	// 自定义版本安装目录,不在统一目录下的版本 key: version, value: install_dir
	LocalVersions map[string]string `json:"local_versions"`
}

// ActiveEnvNames 返回激活环境变量列表
func (t *ToolChain) ActiveEnvNames() []string {
	return maputil.TypedKeys(t.ActiveEnv)
}

// RenderActiveEnv 渲染激活环境变量值中的一些表达式变量 eg: {version}, {install_dir}
func (t *ToolChain) RenderActiveEnv(varMap map[string]string) map[string]string {
	realEnvMap := make(map[string]string)
	for k, val := range t.ActiveEnv {
		if strings.Contains(val, "{") {
			val = strutil.ReplaceVars(val, varMap)
		}
		realEnvMap[k] = val
	}
	return realEnvMap
}

// SimpleTool 简单独立工具 - 单文件，可执行，不需要多版本处理的工具，只需安装最新的即可。PortableTool, StaticTool
//   - 例如 `curl`, `wget`, `ast-grep`, `ripgrep` 等工具。
//   - 支持直接从 github 快速下载安装 `xenv tools install --uri github:user/repo rg@latest`
//   - 支持从任意 URL 下载安装 `xenv tools install --uri https://example.com/file.tar.gz`
type SimpleTool struct {
	// ID         string `json:"id"`
	Name       string `json:"name"`
	InstallURL string `json:"install_url"`
	InstallDir string `json:"install_dir"`
	BinName    string `json:"bin_name"`
	Version    string `json:"version"` // 版本号，如 "1.21", "lts", "latest"
}

// VersionSpec 版本规格
type VersionSpec struct {
	Name    string // SDK名称
	Version string // 输入的版本 可以是 lts, latest
	// 实际的版本号 eg: 1.21.1
	RealVersion string
	// Global bool    // scope: global
}

// ID 返回版本规格的ID name:version
func (vs *VersionSpec) ID() string { return vs.String() }

// RealID 返回版本规格的ID name:real_version
func (vs *VersionSpec) RealID() string { return vs.Name + ":" + vs.RealVersion }

// String 返回版本规格的字符串表示
func (vs *VersionSpec) String() string {
	return vs.Name + ":" + vs.Version
}
