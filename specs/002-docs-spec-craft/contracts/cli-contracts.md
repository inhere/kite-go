# CLI Command Contracts: Kite XEnv

## Command Structure
所有xenv命令遵循以下结构：
`kite xenv <subcommand> [options] [arguments]`

## Subcommands Contracts

### 1. xenv config
**Purpose**: 配置管理和查看

#### Actions:
- **View Config** (Default): 显示当前配置信息
  - Input: 无特定参数
  - Output: 显示当前使用的配置文件路径和配置详情
  - Format: human-readable 或 JSON

- **Set Config**: 设置配置项
  - Input: `--set <name> <value>`
  - Output: 配置项设置成功确认
  - Error: 无效配置名称或值

- **Export Config**: 导出配置
  - Input: `--export [zip|json]`
  - Output: 导出文件路径
  - Error: 文件写入失败或磁盘空间不足

- **Import Config**: 从文件导入配置
  - Input: `--import <path>`
  - Output: 导入成功确认
  - Error: 文件不存在或格式错误

### 2. xenv list
**Purpose**: 列出所有已安装的工具链和环境设置信息

#### Actions:
- **List Tools** (Default, with `--tool`): 列出已安装的工具链
  - Input: `[--tool]`
  - Output: 工具链列表，标记已激活的版本
  - Format: human-readable 或 JSON

- **List Env** (with `--env`): 列出环境变量
  - Input: `--env`
  - Output: 环境变量列表
  - Format: human-readable 或 JSON

- **List Path** (with `--path`): 列出PATH路径
  - Input: `--path`
  - Output: PATH路径列表
  - Format: human-readable 或 JSON

- **List Activity** (with `--activity`): 列出已激活的项目
  - Input: `--activity`
  - Output: 已激活的工具链和路径
  - Format: human-readable 或 JSON

- **List All** (with `--all`): 列出所有设置
  - Input: `--all`
  - Output: 所有已安装的工具链和环境设置信息
  - Format: human-readable 或 JSON

### 3. xenv init
**Purpose**: 导入配置后执行初始化处理

#### Actions:
- **Initialize Config**: 激活配置中指定的工具链
  - Input: 无特定参数
  - Output: 初始化进度和结果
  - Error: 配置文件错误或工具安装失败

### 4. xenv shell
**Purpose**: shell集成

#### Actions:
- **Generate Shell Hook**: 生成shell集成脚本
  - Input: `--type [pwsh|bash|zsh]`
  - Output: shell命令脚本
  - Error: 未知shell类型

### 5. xenv env
**Purpose**: 环境变量管理

#### Actions:
- **List Environment Variables** (Default): 列出所有环境变量
  - Input: 无特定参数
  - Output: 环境变量列表
  - Format: human-readable 或 JSON

- **Set Environment Variable**: 设置环境变量
  - Input: `--set [-g] <name> <value>`
  - Output: 变量设置成功确认
  - Error: 无效变量名或值

- **Unset Environment Variable**: 删除环境变量
  - Input: `--unset [-g] <name...>`
  - Output: 变量删除确认
  - Error: 变量不存在

### 6. xenv path
**Purpose**: PATH管理

#### Actions:
- **List Paths** (Default): 列出所有PATH路径
  - Input: 无特定参数
  - Output: PATH路径列表
  - Format: human-readable 或 JSON

- **Add Path**: 添加PATH路径
  - Input: `--add [-g] <path>`
  - Output: 路径添加成功确认
  - Error: 路径不存在或已存在

- **Remove Path**: 删除PATH路径
  - Input: `--rm [-g] <path>`
  - Output: 路径删除成功确认
  - Error: 路径不存在

- **Search Path**: 搜索PATH中的路径
  - Input: `-s <path>`
  - Output: 匹配的路径列表
  - Error: 无匹配结果

### 7. xenv tools
**Purpose**: 工具链管理

#### Actions:
- **List Tools** (Default): 列出管理的工具链
  - Input: 无特定参数
  - Output: 工具链列表，标记当前激活版本
  - Format: human-readable 或 JSON

- **Install Tool**: 安装工具链
  - Input: `install <name:version> ...`
  - Output: 安装进度和成功确认
  - Error: 工具不存在或下载失败

- **Uninstall Tool**: 卸载工具链
  - Input: `uninstall <name:version> ...`
  - Output: 卸载成功确认和选项保留配置
  - Error: 工具未安装

- **Update Tool**: 更新工具链
  - Input: `update <name:version> ...`
  - Output: 更新进度和成功确认
  - Error: 工具未安装或更新失败

- **Show Tool**: 显示工具信息
  - Input: `show <name>`
  - Output: 工具详细信息
  - Error: 工具不存在

### 8. xenv use
**Purpose**: 切换当前环境工具链

#### Actions:
- **Activate Toolchain**: 激活指定工具链版本
  - Input: `[-g] <name:version> ...`
  - Output: 激活确认
  - Error: 工具未安装或版本不存在

### 9. xenv unuse
**Purpose**: 去掉当前环境工具链

#### Actions:
- **Deactivate Toolchain**: 停用指定工具链版本
  - Input: `[-g] <name:version> ...`
  - Output: 停用确认
  - Error: 工具未激活或不存在

## Common Options
- `-g, --global`: 应用到全局而非当前会话
- `-h, --help`: 显示帮助信息
- `--json`: 以JSON格式输出结果

## Error Codes
- 0: 成功
- 1: 通用错误
- 2: 无效参数
- 3: 资源未找到
- 4: 权限错误
- 5: 网络错误
- 6: 配置错误