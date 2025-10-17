# Data Model: Kite XEnv

## Entities

### ToolChain
- **ID**: string (唯一标识符，格式为 name:version)
- **Name**: string (工具名称，如 "go", "node")
- **Version**: string (版本号，如 "1.21", "lts", "latest")
- **Alias**: []string (别名列表，如 ["golang"] for go)
- **InstallURL**: string (可选，下载URL模板)
- **InstallDir**: string (安装目录路径)
- **ActiveEnv**: map[string]string (激活时设置的额外环境变量)
- **Installed**: bool (是否已安装)
- **BinPaths**: []string (该工具的二进制文件路径列表)

### EnvironmentVariable
- **Name**: string (环境变量名称)
- **Value**: string (环境变量值)
- **Scope**: string (作用域: "global" 或 "session")
- **IsActive**: bool (是否当前激活)

### PathEntry
- **Path**: string (添加到PATH的路径)
- **Priority**: int (优先级，数值越小优先级越高)
- **Scope**: string (作用域: "global" 或 "session")
- **IsActive**: bool (是否当前激活)

### Configuration
- **BinDir**: string (默认: ~/.local/bin)
- **InstallDir**: string (默认: ~/.xenv/tools)
- **ShellScriptsDir**: string (默认: ~/.config/xenv/hooks/)
- **Tools**: []ToolChain (可管理的工具链列表)
- **GlobalEnv**: map[string]EnvironmentVariable (全局环境变量)
- **GlobalPaths**: []PathEntry (全局PATH条目)

### ActivityState
- **ActiveTools**: map[string]string (激活的工具链映射，key为工具名，value为版本)
- **ActiveEnv**: map[string]string (激活的环境变量)
- **ActivePaths**: []string (激活的路径列表)
- **LastUpdated**: time.Time (最后更新时间)

### User
- **ID**: string (用户唯一标识符)
- **ConfigPath**: string (用户配置文件路径)
- **HomeDir**: string (用户主目录)
- **ShellType**: string (用户使用的shell类型: "bash", "zsh", "pwsh")