package models

import "time"

// Configuration 代表用户的配置信息，包含工具管理设置、路径配置、环境激活状态等数据
type Configuration struct {
	// ID              string                  `json:"id"`
	// tools 可执行文件链接目录 默认: ~/.local/bin
	BinDir string `json:"bin_dir"`
	// 工具安装基础目录 默认: ~/.xenv/tools
	InstallDir string `json:"install_dir"`
	// shell hooks 脚本目录。 默认: ~/.config/xenv/hooks/
	ShellHooksDir string `json:"shell_hooks_dir"`
	// 全局环境变量 - 首次初始化生效，后续通过命令设置即可
	GlobalEnv map[string]EnvVariable `json:"global_env"`
	// 全局PATH条目 - 首次初始化生效，后续通过命令设置即可
	GlobalPaths []PathEntry `json:"global_paths"`
	// 从远程下载不同OS平台的工具包的后缀格式
	// eg:
	//
	// 	os_download_ext:
	// 	  windows: zip
	// 	  linux: tar.gz
	// 	  macos: tar.gz
	OSDownloadExt map[string]string `json:"os_download_ext"`
	DownloadDir   string            `json:"download_dir"` // 临时下载目录
	Tools         []ToolChain       `json:"tools"`        // 可管理的工具链列表
	SimpleTools   []SimpleTool      `json:"simple_tools"` // 配置的简单工具列表
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}
