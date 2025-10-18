package models

// ToolChain SDK开发工具（如Go、Node.js等）配置，包含安装路径、别名等属性。
//   - 只是工具信息配置，不含有特定的版本信息
type ToolChain struct {
	Name  string `json:"name"`  // 工具名称，如 "go", "node"
	Alias string `json:"alias"` // 工具别名列表，如 "golang" for go
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
	ActiveEnv map[string]string `json:"active_env"`
	// 该工具的 bin 文件目录名称，不设置就是 install_dir 目录
	BinDir      string   `json:"bin_dir"`
	BinPaths    []string `json:"bin_paths"`    // 该工具提供的二进制文件路径列表
	PostInstall []string `json:"post_install"` // 安装完成后执行的shell hook脚本
	// 自定义版本安装目录,不在统一目录下的版本 key: version, value: install_dir
	LocalVersions map[string]string `json:"local_versions"`
}

// SimpleTool 简单独立工具 - 单文件，可执行，不需要多版本处理的工具，只需安装最新的即可。PortableTool, StaticTool
//   - 例如 `curl`, `wget`, `ast-grep`, `ripgrep` 等工具。
//   - 支持直接从 github 快速下载安装 `xenv tools install --uri github:user/repo rg@latest`
//   - 支持从任意 URL 下载安装 `xenv tools install --uri https://example.com/file.tar.gz`
type SimpleTool struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	InstallURL string `json:"install_url"`
	InstallDir string `json:"install_dir"`
	BinName    string `json:"bin_name"`
	Version    string `json:"version"` // 版本号，如 "1.21", "lts", "latest"
}
