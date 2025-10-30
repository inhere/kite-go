package models

// Configuration 代表用户的配置信息，包含工具管理设置、路径配置、环境激活状态等数据
// TODO 将外部工具独立出去， xenv 只管理 ENV, PATH, SDK
type Configuration struct {
	// ID              string                  `json:"id"`
	// tools 可执行文件链接目录 默认: ~/.local/bin
	BinDir string `json:"bin_dir"`
	// 工具安装基础目录 默认: ~/.xenv/tools
	InstallDir string `json:"install_dir"`
	// 快速配置 shell 命令别名, 会自动注入到shell环境
	ShellAliases map[string]string `json:"shell_aliases"`
	// shell hooks 脚本目录。 默认: ~/.config/xenv/hooks/
	ShellHooksDir string `json:"shell_hooks_dir"`
	// 全局环境变量 - 首次初始化生效，后续通过命令设置即可
	GlobalEnv map[string]string `json:"global_env"`
	// 全局PATH条目 - 首次初始化生效，后续通过命令设置即可
	GlobalPaths []string `json:"global_paths"`
	// 从远程下载不同OS平台的工具包的后缀格式
	// eg:
	//
	// 	download_ext:
	// 	  windows: zip
	// 	  linux: tar.gz
	// 	  darwin: tar.gz
	DownloadExt map[string]string `json:"download_ext"`
	DownloadDir   string            `json:"download_dir"` // 临时下载目录
	// 可管理的工具链列表
	//  - sdks 和 tools 差异是：sdk 允许本地同时存在多个版本，tools 只允许一个版本
	SDKs []ToolChain `json:"sdks"`
	// 配置的简单工具列表
	Tools []SimpleTool `json:"tools"`
	// internal fields
	configFile string
	configDir  string
}

// IsDefinedSDK returns true if the SDK configuration is defined
func (c *Configuration) IsDefinedSDK(name string) bool {
	// Check if the tool is installed
	toolFound := false
	for _, tool := range c.SDKs {
		if tool.Name == name {
			toolFound = true
			break
		}
	}
	return toolFound
}

// FindSDKConfig returns the SDK configuration if it is defined
func (c *Configuration) FindSDKConfig(name string) *ToolChain {
	for _, tool := range c.SDKs {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

func (c *Configuration) ConfigFile() string {
	return c.configFile
}

// SetConfigFile sets the config.yaml configuration file path
func (c *Configuration) SetConfigFile(filePath string) {
	c.configFile = filePath
}

// ConfigDir returns the directory path where the configuration file is located
func (c *Configuration) ConfigDir() string {
	return c.configDir
}

// SetConfigDir sets the directory path where the configuration file is located
func (c *Configuration) SetConfigDir(dirPath string) {
	c.configDir = dirPath
}
