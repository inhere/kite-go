# 环境管理工具 (Environment Manager)

Kite CLI 应用的环境管理工具，用于管理本地开发环境中的 SDK 版本切换。

类似工具：

- asdf
- https://github.com/jdx/mise

## 功能特性

- 支持多种开发环境：Go、Node.js、Java、Flutter
- 跨平台 shell 环境集成（bash、zsh、pwsh、cmd）
- 自动 SDK 下载和安装
- 环境隔离和版本管理
- Shell 脚本自动生成

## 快速开始

```bash
echo 'eval "$(kite dev env shell bash)"' >> ~/.bashrc
echo 'eval "$(kite dev env shell zsh)"' >> ~/.zshrc
echo 'kite dev env shell fish | source' >> ~/.config/fish/config.fish
echo 'kite dev env shell pwsh | Out-String | Invoke-Expression' >> ~/.config/powershell/Microsoft.PowerShell_profile.ps1
```

### 1. 启用 Shell 集成

#### Bash
```bash
# 添加到 ~/.bashrc
eval "$(kite dev env shell bash)"
```

#### Zsh  
```bash
# 添加到 ~/.zshrc
eval "$(kite dev env shell zsh)"
```

#### PowerShell
```powershell
# 添加到 $PROFILE
Invoke-Expression (kite dev env shell pwsh)
```

#### CMD
```batch
REM 运行命令生成批处理文件
kite dev env shell cmd > ktenv.bat && call ktenv.bat
```

### 2. 使用 ktenv 命令

#### 安装 SDK
```bash
# 安装特定版本
ktenv add go:1.21.5
ktenv add node:18.17.0

# 安装最新版本
ktenv add go:latest
ktenv add node:lts
```

#### 激活 SDK
```bash
# 激活单个 SDK
ktenv use go:1.21.5

# 激活多个 SDK
ktenv use go:1.21.5 node:18

# 激活并保存到项目配置
ktenv use -s go:1.21.5 node:18
```

#### 列出已安装的 SDK
```bash
# 列出所有 SDK
ktenv list

# 列出特定类型 SDK
ktenv list go
```

#### 取消激活 SDK
```bash
# 取消激活特定 SDK
ktenv unuse go

# 取消激活多个 SDK
ktenv unuse go node
```

## 命令参考

### kite dev env 命令组

| 命令 | 描述 | 参数 |
|------|------|------|
| `kite dev env list` | 列出已安装的SDK | `--type` 过滤SDK类型 |
| `kite dev env add` | 安装新的SDK | `<sdk:version>` |
| `kite dev env use` | 激活SDK版本 | `<sdk:version>` `--save` |
| `kite dev env remove` | 移除已安装的SDK | `<sdk:version>` |
| `kite dev env shell` | 生成shell注入脚本 | `[shell-type]` |
| `kite dev env config` | 查看和编辑配置 | `--edit` |

### ktenv 函数命令

| 命令 | 描述 | 语法 |
|------|------|------|
| `ktenv use` | 激活SDK版本 | `ktenv use <sdk:version>...` |
| `ktenv use -s` | 激活并保存配置 | `ktenv use -s <sdk:version>...` |
| `ktenv unuse` | 取消激活SDK | `ktenv unuse <sdk>...` |
| `ktenv add` | 下载安装SDK | `ktenv add <sdk:version>...` |
| `ktenv list` | 显示SDK状态 | `ktenv list [sdk]` |

## 版本格式

| 格式 | 描述 | 示例 |
|------|------|------|
| 精确版本 | 完整版本号 | `go:1.21.5` |
| 主版本 | 主版本最新 | `node:18` |
| 别名版本 | 预定义别名 | `node:lts`, `go:latest` |
| 自动检测 | 基于项目配置 | `go:auto` |

## 支持的 SDK

### Go
- 下载源：golang.org
- 环境变量：自动设置 GOROOT
- 二进制路径：`{install_dir}/bin`

### Node.js
- 下载源：nodejs.org
- 环境变量：NODE_ENV 等
- 二进制路径：`{install_dir}/bin`

### Java
- 下载源：Oracle JDK
- 环境变量：JAVA_HOME
- 二进制路径：`{install_dir}/bin`

### Flutter
- 手动安装到指定目录
- 环境变量：FLUTTER_HOME
- 二进制路径：`{install_dir}/bin`

## 配置文件

### 主配置文件
位置：`~/.kite-go/config/module/shell_env.yml`

```yaml
add_paths: []              # 全局PATH路径
add_envs: {}              # 全局环境变量
remove_envs: []           # 需要移除的环境变量
sdk_dir: /opt/devsdk      # SDK安装基础目录

sdks:
  - name: go
    install_url: https://golang.org/dl/go{version}.{os}-{arch}.tar.gz
    install_dir: /opt/devsdk/go{version}
    
  - name: node
    install_url: https://nodejs.org/dist/v{version}/{os}-{arch}.tar.xz
    install_dir: /opt/devsdk/node{version}
    
  - name: java
    install_url: https://download.oracle.com/otn_software/java/jdk/{version}/jdk-{version}_linux-{arch}.tar.gz
    install_dir: /opt/devsdk/java{version}
    active_env:
      JAVA_HOME: /opt/devsdk/java{version}
```

### 活跃状态文件
位置：`~/.kite-go/data/shell_env/active.json`

```json
{
  "current_sdks": {
    "go": "1.21.5",
    "node": "18.17.0"
  },
  "add_paths": [
    "/opt/devsdk/go1.21.5/bin",
    "/opt/devsdk/node18.17.0/bin"
  ],
  "add_envs": {
    "GOROOT": "/opt/devsdk/go1.21.5",
    "JAVA_HOME": "/opt/devsdk/java11"
  },
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 自定义脚本

可以在 `~/.kite-go/data/shell_env/` 目录下创建自定义脚本：

- `init.sh` - 初始化脚本
- `hooks/` - 钩子脚本目录
- `functions/` - 自定义函数目录

## 故障排除

### 1. ktenv 命令不存在
确保已正确添加 shell 集成代码到配置文件中。

### 2. SDK 下载失败
检查网络连接和下载 URL 配置。

### 3. 权限问题
确保对 SDK 安装目录有写权限。

### 4. 路径未生效
重新启动 shell 会话或运行 `source ~/.bashrc`。

## 开发者指南

### 架构概览
```
pkg/envmgr/
├── types.go         # 核心数据结构
├── version.go       # 版本解析
├── config.go        # 配置管理
├── state.go         # 状态管理
├── sdk.go           # SDK管理
├── shell.go         # Shell脚本生成
└── manager.go       # 主管理器
```

### 添加新的 SDK 支持

1. 在配置文件中添加 SDK 定义
2. 实现特定的环境变量和路径逻辑
3. 添加下载和安装逻辑
4. 更新文档和测试

### 扩展 Shell 支持

1. 在 `ShellType` 枚举中添加新类型
2. 在脚本生成器中实现对应逻辑
3. 添加测试用例

## 相关项目

- https://github.com/direnv/direnv - a tool for managing your dir environment
* [autoenv](https://github.com/hyperupcall/autoenv) - older, popular, and lightweight.
* [zsh-autoenv](https://github.com/Tarrasch/zsh-autoenv) - a feature-rich mixture of autoenv and [smartcd](https://github.com/cxreg/smartcd): enter/leave events, nesting, stashing (Zsh-only).
* [asdf](https://github.com/asdf-vm/asdf) - a pure bash solution that has a plugin system. The [asdf-direnv](https://github.com/asdf-community/asdf-direnv) plugin allows using asdf managed tools with direnv.
* [ondir](https://github.com/alecthomas/ondir) - OnDir is a small program to automate tasks specific to certain directories
* [shadowenv](https://shopify.github.io/shadowenv/) - uses an s-expression format to define environment changes that should be executed
* [quickenv](https://github.com/untitaker/quickenv) - an alternative loader for `.envrc` files that does not hook into your shell and favors speed over convenience.
* [mise](https://github.com/jdx/mise) - direnv, make and asdf all in one tool.
