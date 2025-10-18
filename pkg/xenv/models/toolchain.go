package models

// ToolChain SDK开发工具（如Go、Node.js等）配置，包含安装路径、别名等属性。
//  - 只是工具信息配置，不含有特定的版本信息
type ToolChain struct {
	ID        string   `json:"id"`        // 唯一标识符，格式为 name:version
	Name      string   `json:"name"`      // 工具名称，如 "go", "node"
	Alias     []string `json:"alias"`     // 别名列表，如 ["golang"] for go
	Installed bool     `json:"installed"` // 是否已安装 - 本地至少有一个版本的工具
	// 可选，下载URL模板 eg: "https://golang.org/dl/go{version}.{os}-{arch}.{download_ext}"
	InstallURL  string            `json:"install_url"`
	// sdk tool 安装目录路径 默认 ~/.xenv/tools/{name}/{version}
	InstallDir  string `json:"install_dir"`
	ActiveEnv   map[string]string `json:"active_env"`   // 激活时设置的额外环境变量
	InitVersion string `json:"init_version"`            // 首次初始化时确认下载的 tool 版本号
	BinDir      string            `json:"bin_dir"`      // 该工具的 bin 文件目录名称，不设置就是 install_dir 目录
	BinPaths    []string          `json:"bin_paths"`    // 该工具提供的二进制文件路径列表
	PostInstall []string          `json:"post_install"` // 安装完成后执行的shell hook脚本
}

// SimpleTool 简单独立工具 - 单文件，可执行，不需要多版本处理的工具，只需安装最新的即可。PortableTool, StaticTool
//   - 例如 `curl`, `wget`, `ast-grep`, `ripgrep` 等工具。
//   - 支持直接从 github 快速下载安装 `xenv tools install --uri github:user/repo rg@latest`
//   - 支持从任意 URL 下载安装 `xenv tools install --uri https://example.com/file.tar.gz`
type SimpleTool struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InstallURL string    `json:"install_url"`
	InstallDir string    `json:"install_dir"`
	BinName    string    `json:"bin_name"`
	Version string `json:"version"` // 版本号，如 "1.21", "lts", "latest"
}
