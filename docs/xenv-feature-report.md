# kite xenv 功能报告

## 1. 模块定位

`kite xenv` 是 Kite CLI 中用于管理本机开发环境的功能模块，目标类似 `mise`、`vfox`、`asdf`：

- 管理本地 SDK/工具链的多版本安装与激活
- 管理环境变量和 `PATH`
- 支持全局、当前 shell 会话、项目目录三种作用域
- 通过 shell hook 让环境变更立即作用于当前终端
- 支持通过 `.xenv.toml` 实现目录级开发环境配置

当前模块入口位于：

- CLI 命令：`internal/cli/xenvcmd`
- 核心能力：`pkg/xenv`
- 默认配置样例：`data/xenv/config.yaml`
- 功能需求草稿：`docs/feat-craft/kite-xenv-spec-craft.md`

## 2. 功能总览

### 2.1 SDK/工具链管理

`xenv` 支持管理 Go、Node、Flutter、Python 等可多版本共存的 SDK 工具链。

常用命令：

```bash
kite xenv tools list
kite xenv tools index
kite xenv tools install go:1.22.0
kite xenv tools uninstall go:1.22.0
kite xenv tools update go:1.22.0
kite xenv tools show go
kite xenv use go:1.22
kite xenv unuse go:1.22
```

版本规格支持：

```bash
kite xenv use go
kite xenv use go:1.22
kite xenv use go@1.22
kite xenv use go:latest
```

说明：

- `use` / `unuse` 通过 `tools.ParseVersionSpec()` 解析版本，支持 `name`、`name:version`、`name@version`。
- `tools install`、`tools uninstall`、`tools update` 当前 CLI 参数解析更严格，要求 `name:version` 格式。
- 本地已安装 SDK 元数据保存到 `~/.xenv/tools.local.json`。

### 2.2 环境变量管理

支持设置、取消和查看环境变量。

```bash
kite xenv env list
kite xenv env set FOO bar
kite xenv env unset FOO
```

快捷写法：

```bash
kite xenv set FOO bar
kite xenv unset FOO
```

作用域参数：

```bash
# 当前 shell 会话
kite xenv env set FOO bar

# 全局状态
kite xenv env set -g FOO bar

# 当前目录 .xenv.toml
kite xenv env set -s FOO bar
kite xenv env set -d FOO bar
```

环境变量名称会被转换为大写，并校验是否为合法变量名。

### 2.3 PATH 管理

支持添加、删除、搜索、查看 `PATH` 条目。

```bash
kite xenv path list
kite xenv path add ./bin
kite xenv path remove ./bin
kite xenv path search go
```

作用域参数：

```bash
kite xenv path add -g ~/.local/bin
kite xenv path add -s ./bin
```

说明：

- `path add` 会先规范化路径，并检查目录是否存在。
- `path remove` 当前会检查路径是否存在于当前进程的 `PATH` 中。
- 添加路径时会放到 `PATH` 前部，使其拥有更高优先级。

### 2.4 Shell Hook 集成

`xenv` 的核心使用方式依赖 shell hook。未配置 shell hook 时，`env`、`path`、`use` 命令可以更新状态文件，但不能直接修改当前终端环境。

生成 hook：

```bash
kite xenv shell --type bash
kite xenv shell --type zsh
kite xenv shell --type pwsh
```

Bash：

```bash
eval "$(kite xenv shell --type bash)"
```

Zsh：

```bash
eval "$(kite xenv shell --type zsh)"
```

PowerShell：

```powershell
Invoke-Expression (& kite xenv shell --type pwsh)
```

或者：

```powershell
kite xenv shell --type pwsh | Out-String | Invoke-Expression
```

hook 会设置：

```text
XENV_HOOK_SHELL
XENV_SESSION_ID
```

配置 hook 后，推荐使用注入到 shell 的 `xenv` 函数：

```bash
xenv use go:1.22
xenv set FOO bar
xenv path add ./bin
```

`xenv` 函数会执行 `kite xenv ...`，并解析输出中的 `--Expression--` 标记，将后半部分作为 shell 脚本在当前终端执行。

### 2.5 目录级环境

`xenv` 支持当前目录或父目录中的 `.xenv.toml`，用于项目级环境配置。

示例：

```toml
paths = [
  "./bin",
]

[sdks]
go = "1.22"
node = "20"

[envs]
APP_ENV = "local"
DEBUG = "true"

[tools]
ripgrep = "*"
```

保存到目录配置：

```bash
xenv use -s go:1.22
xenv set -s APP_ENV local
xenv path add -s ./bin
```

hook 会重写或包装 `cd` / `Set-Location`，进入目录后调用内部命令：

```bash
kite xenv init-direnv
```

然后根据最近的 `.xenv.toml` 激活 SDK、环境变量和路径。

当前实现只查找最近的一个 `.xenv.toml`，尚未实现多层目录状态叠加。

## 3. 配置文件

默认配置文件：

```text
~/.config/xenv/config.yaml
```

默认配置目录：

```text
~/.config/xenv/
```

默认状态和工具索引：

```text
~/.config/xenv/global.toml
~/.xenv/session/<session_id>.json
~/.xenv/tools.local.json
```

默认配置项：

```yaml
bin_dir: "~/.local/bin"
install_dir: "~/.xenv/tools"
shell_hooks_dir: "~/.config/xenv/hooks/"
global_env: {}
global_paths: []
shell_aliases: {}
download_ext:
  windows: zip
  linux: tar.gz
  darwin: tar.gz
sdks: []
tools: []
```

SDK 配置示例：

```yaml
sdks:
  - name: go
    alias: golang
    install_url: "https://golang.org/dl/go{version}.{os}-{arch}.{download_ext}"
    install_dir: "D:/work/env/devsdk/gosdk/go{version}"
    active_env:
      GO111MODULE: auto
    bin_dir: bin
    other_versions:
      latest: "C:/Users/inhere/scoop/apps/go/current"

  - name: node
    install_url: "https://cdn.npmmirror.com/binaries/node/v{version}/node-v{version}-{os}-{arch}.{download_ext}"
    install_dir: "D:/work/env/devsdk/nodejs/node-v{version}-win-x64"
    download_ext:
      windows: zip
```

注意：当前模型字段是 `other_versions`。`data/xenv/config.yaml` 中存在 `local_versions` 示例字段，按当前代码应优先使用 `other_versions`。

## 4. 状态模型

`xenv` 的激活状态分三层：

```text
全局状态: ~/.config/xenv/global.toml
目录状态: 当前目录或父目录的 .xenv.toml
会话状态: ~/.xenv/session/<session_id>.json
```

加载顺序：

```text
global -> direnv -> session
```

合并后的状态用于生成 shell 初始化脚本和激活 SDK。

状态内容主要包含：

```toml
paths = []

[sdks]
go = "1.22"

[envs]
APP_ENV = "local"

[tools]
ripgrep = "*"
```

查看状态：

```bash
kite xenv list activity
kite xenv list activity -t
```

## 5. 初始化和推荐使用流程

### 5.1 初始化

```bash
kite xenv init
```

该命令会：

- 加载或创建 `~/.config/xenv/config.yaml`
- 创建 `bin_dir`
- 创建 `install_dir`
- 创建 `shell_hooks_dir`

### 5.2 手动安装 SDK 后纳入管理

当前 SDK 自动下载安装仍不完整，推荐先手动安装 SDK，再让 `xenv` 索引。

配置 `~/.config/xenv/config.yaml`：

```yaml
sdks:
  - name: go
    install_dir: "D:/work/env/devsdk/gosdk/go{version}"
    active_env:
      GO111MODULE: auto
    bin_dir: bin
    other_versions:
      latest: "C:/Users/inhere/scoop/apps/go/current"
```

索引本地工具：

```bash
kite xenv tools index
```

查看工具：

```bash
kite xenv tools list
```

激活版本：

```bash
xenv use go:1.22
```

设置全局默认：

```bash
xenv use -g go:1.22
```

设置项目目录专用：

```bash
xenv use -s go:1.22
```

## 6. 命令清单

### 6.1 主命令

```bash
kite xenv
```

子命令：

```text
tools
use
unuse
env
path
config
list
init
shell
shell-init-hook
shell-direnv
```

隐藏内部命令：

```text
shell-init-hook
shell-direnv
```

### 6.2 tools

```bash
kite xenv tools install <name:version>...
kite xenv tools uninstall <name:version>
kite xenv tools update <name:version>...
kite xenv tools show <name>
kite xenv tools list
kite xenv tools register
kite xenv tools index
```

别名：

```text
tools: t, tool, sdks, sdk
install: i, in
uninstall: un, rm, remove
update: up
list: ls
index: idx
register: add, reg
```

### 6.3 use / unuse

```bash
kite xenv use [-g] [-s|-d] <name:version>...
kite xenv unuse [-g] [-s|-d] <name:version>...
```

作用域：

```text
默认: 当前会话
-g: 全局
-s, -d: 当前目录 .xenv.toml
```

### 6.4 env

```bash
kite xenv env
kite xenv env list
kite xenv env set [-g] [-s|-d] <name> <value>
kite xenv env unset [-g] [-s|-d] <name...>
```

快捷命令：

```bash
kite xenv set <name> <value>
kite xenv unset <name...>
```

### 6.5 path

```bash
kite xenv path
kite xenv path list
kite xenv path add [-g] [-s|-d] <path>
kite xenv path remove [-g] [-s|-d] <path>
kite xenv path search <value>
```

别名：

```text
remove: rm, delete
list: ls
search: s
```

### 6.6 list

```bash
kite xenv list
kite xenv list tools
kite xenv list env
kite xenv list path
kite xenv list activity
kite xenv list all
```

说明：

- `kite xenv list` 默认列出工具链。
- `list all` 当前仍是占位实现。

### 6.7 config

```bash
kite xenv config
kite xenv config get <name>
kite xenv config set <name> <value>
kite xenv config export [zip|json]
kite xenv config import <path>
```

支持的 `get` / `set` 配置项：

```text
bin_dir
install_dir
shell_hooks_dir
```

注意：当前 `config set` 和 `config import` 的保存写回逻辑仍是 TODO，执行后需要核对配置文件是否实际更新。

## 7. 当前实现完成度

可用度较高：

- `xenv init`
- shell hook 脚本生成
- 当前 shell 中的 `xenv` 包装函数
- 环境变量设置、删除、查看
- `PATH` 添加、删除、搜索、查看
- `.xenv.toml` 目录状态加载和保存
- 本地 SDK 索引 `tools index`
- SDK 列表 `tools list`
- 激活本地已索引 SDK `use`

仍需完善：

- `shell --install` 调用的 `InstallToProfile()` 当前基本为空实现，建议手动写入 shell profile。
- `tools install` 的 zip/tar.gz 解压函数仍是 TODO。
- `config set`、`config import` 的保存写回逻辑仍是 TODO。
- `tools register` 当前返回 TODO。
- `list all` 当前是占位输出。
- `unuse` 对状态中 SDK key 的删除逻辑需要进一步验证。
- `.envrc` / `.envrc.ps1` 目前有发现逻辑，但实际执行集成仍不完整。

## 8. 最小可用示例

独立 CLI：

```bash
go run ./cmd/xenv init
go run ./cmd/xenv tools list
go run ./cmd/xenv shell --type bash
```

构建独立二进制：

```bash
make build-xenv
make install-xenv
```

PowerShell：

```powershell
kite xenv init
Invoke-Expression (& kite xenv shell --type pwsh)

kite xenv tools index
kite xenv tools list

xenv use go:latest
xenv set DEMO_ENV hello
xenv path add ./bin

kite xenv list activity
```

Bash：

```bash
kite xenv init
eval "$(kite xenv shell --type bash)"

kite xenv tools index
kite xenv tools list

xenv use go:latest
xenv set DEMO_ENV hello
xenv path add ./bin

kite xenv list activity
```

安装独立 `xenv` 二进制后，可以直接使用：

```bash
xenv init
eval "$(xenv shell --type bash)"

xenv tools index
xenv tools list
xenv use go:latest
```

项目目录专用环境：

```bash
xenv use -s go:1.22
xenv set -s APP_ENV local
xenv path add -s ./bin
```

执行后会生成或更新当前目录的 `.xenv.toml`。

## 9. 使用建议

现阶段推荐将 `xenv` 用作“本地已安装 SDK 的激活器”和“项目环境变量/PATH 管理器”：

- SDK 先通过系统包管理器或手动方式安装。
- 在 `config.yaml` 中声明 SDK 安装目录。
- 使用 `kite xenv tools index` 建立本地索引。
- 使用 shell hook 中的 `xenv use` 激活版本。
- 使用 `.xenv.toml` 固化项目级环境。

自动下载安装、配置导入导出和 profile 自动写入能力已经有结构和命令入口，但仍需要补齐实现后再作为主流程使用。
