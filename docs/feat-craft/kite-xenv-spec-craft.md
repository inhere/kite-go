# kite xenv 命令需求草稿

为 kite cli新增实现一个命令 `xenv` 用于管理本机开发环境和工具，类似 `mise` `vfox` 工具。
主要分两块功能：tool 管理，env/path 管理。
同时希望实现跨平台支持，在另一台机器上导入配置即可以使用。

## 配置说明

- user 配置文件保存在 `~/.config/xenv/config.yaml`
- `shell_hooks_dir` 配置项，用于在 `xenv shell` 注入hooks中执行自定义脚本文件，默认为 `~/.config/xenv/hooks/`
- `bin_dir` 配置项，用于指定工具安装后创建的软链接/shims目录，默认为 `~/.local/bin`

## 基础功能

- `kite xenv config` xenv 配置管理和查看等
  - 默认显示配置信息，列出使用的配置文件
  - `--set <name> <value>` 设置配置项. name 支持路径方式 eg: `xenv config --set xxx.yyy.zzz value`
  - `--export [zip|json]` 导出配置，将相关的配置文件打包到 zip 文件或者单个json文件中。
  - `--import <path>` 从文件、URL导入配置
- `kite xenv list` 列出所有已安装的工具链和环境设置信息
  - `--tool` 列出所有已安装的工具链(默认)
  - `--env` 列出所有已设置的环境变量
  - `--path` 列出所有已设置的PATH路径
  - `--activity` 列出已经激活的工具链和路径
  - `--all` 列出所有已安装的工具链和环境设置信息
- `kite xenv init` 导入配置后，执行初始化处理(已经激活的立即检查并安装)

## 核心功能

跨平台的 shell hooks 支持。

- `kite xenv shell --type [pwsh|bash|zsh]` shell 集成，便于 xenv use 后立即生效
  - 将会在 shell 中注入方法 `xenv` 用于让 `xenv env`, `xenv use`, `xenv path` 等命令立即生效
  - 内置shell方法： `use_tool`, `path_add` 方便使用
- 需要在 .bashrc / .zshrc / pwsh profile 中配置hook eg: `eval "$(kite xenv shell --type bash)"`
- hooks 同时会执行 `shell_hooks` 配置的 `~/.config/xenv/hooks/*.{sh|ps1}` 脚本文件
- 内置支持目录下 `.xenv.toml` 文件: 可以配置 tools, env和path路径
- 兼容支持：如果当前目录下存在 `.envrc` 文件，也会自动执行 `.envrc` 文件(仅支持bash,zsh)
- 兼容支持：如果当前目录下存在 `.envrc.ps1` 文件，也会自动执行 `.envrc.ps1` 文件(仅支持pwsh)

## 环境管理

- `kite xenv env` 环境管理。默认列出所有已设置的环境变量和PATH路径
  - `kite xenv env --set [-g] <name> <value>` 添加环境变量
  - `kite xenv env --unset [-g] <name...>` 删除环境变量
  - `-g` 全局有效 将会保存到全局 `~/.config/xenv/activity.json` 中的 `env` object
- `kite xenv path` 环境PATH管理。默认列出所有已设置的PATH路径
  - `kite xenv path --add [-g] <path>` 添加PATH路径
  - `kite xenv path --rm [-g] <path>` 删除PATH路径
  - `-s` 搜索PATH中的路径
  - `-g` 全局有效 将会保存到全局 `~/.config/xenv/activity.json` 中的 `paths` list

## 工具链管理

通过 `kite xenv tools` 命令进行本机开发工具管理。支持多版本开发工具管理下载，也支持独立小工具下载生意。

- `~/.config/xenv/config.yaml` 配置 tools 项，用于配置手动安装的工具链目录和可以通过 http 进行自动安装的工具链。
- `install_dir` string 配置项，用于指定工具链默认安装目录，默认为 `~/.xenv/tools`
- `tools` list 配置项，用于配置手动安装的工具链目录和可以通过 http 进行自动安装的工具链。
  - 在这里配置了的工具，才会被纳入管理。
- 独立工具 - 单文件，可执行，不需要多版本处理的工具，只需安装最新的即可。 
  - 例如 `curl`, `wget`, `ast-grep`, `ripgrep` 等工具。
  - 支持直接从 github 快速下载安装 `xenv tools install --uri github:user/repo rg@latest`

```yaml
bin_dir: ~/.local/bin
install_dir: ~/.xenv/tools

tools:
  - name: go
    alias: [golang]
    install_url: "https://golang.org/dl/go{version}.{os}-{arch}.{isWindows ? zip : tar.gz}" # 可选，不配置则使用本地已下载的工具链
    install_dir: /opt/devsdk/go{version} # 可选，默认使用顶级 `install_dir` 配置项
    active_env: # 可选，激活时额外设置ENV
      GO111MODULE: auto
```

当前版本只支持简单的http方式下载。或者手动安装后，配置文件里配置工具链版本和目录。后续支持添加通过 `asdf` `vfox` `brew` `scoop` 等成熟工具进行下载安装

- `kite xenv tools` 支持多版本的常用开发工具管理。eg: node, go, flutter, php, java, rust 等
  - 默认列出所有管理的工具链，当前激活的会标记出来
- `kite xenv use [-g] <name:version> ...` 切换当前环境工具链
    - eg1: `kite xenv use node@12 go@1.21` 当前会话有效
    - `-g` 会保存到全局 `~/.config/xenv/activity.json` 中的 `tools` object, 重启后仍然有效
- `kite xenv unuse [-g] <name:version> ...` 去掉当前环境工具链
- `kite xenv tools install <name:version> ...` 安装工具链
- `kite xenv tools uninstall <name:version> ...` 卸载工具链
- `kite xenv tools update <name:version> ...` 更新工具链
- `kite xenv tools show <name>` 显示指定工具链信息

> `version` 不设置则使用最新版本 `latest`, 同时支持 `lts` 等特殊标识符
