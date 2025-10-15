# TODO

## kite app

- [ ] generate metadata file for kite
- [ ] 将所有命令和参数生成 json 文件, 用于搜索匹配
- [ ] 使用统一的 app shell 命令生成不同shell环境的脚本
  - 内置支持 bash, zsh, fish, pwsh, cmd(clink) 等
  - 内置支持常用别名设置
  - 支持自定义脚本内容

### kite ext

- [x] 支持注册、管理外部脚本命令
- [x] 支持调用外部命令 eg: 调用 php 实现的命令

### kite plugin

- [ ] 实现类似 git plugin 的事件监听机制
- [ ] go plugin 支持 https://github.com/hashicorp/go-plugin
- [ ] 表达式支持
    - https://github.com/expr-lang/expr Go 的表达式语言和表达式评估
    - https://github.com/hashicorp/go-bexpr Go 中的通用布尔表达式评估
    - https://github.com/ganigeorgiev/fexpr 支持解析类似 SQL 表达式 eg: `id=123 && status='active'`

### run anything

- [x] scripts manage, search, load, parse and run 
- [x] script-task 功能增强
  - [x] 支持扫描工作目录和父级目录的 task 定义文件
  - [x] 运行时支持解析 全局，定义，上下文等变量
  - [x] task command 支持独立定义 workdir, vars 等属性
  - [x] task, command 定义的 vars 支持动态变量 eg: "@sh: git version"
  - [ ] task, command 支持 `If/Cond` 条件表达式
- [ ] 功能增强，支持运行独立的 script-app 文件

## backend serve

- [ ] support backend server: `kited` `kite app:server`
- [ ] cache metadata, provide quick command search

### job/task serve

- [ ] start a background server, listen an port/sock
- [ ] can delivery task by tcp connection

## http tools

- [ ] provide web UI for some operation
- [ ] can delivery task by web page
- [ ] http benchmark tool
- [ ] send http request tool. like curl, ide-http-client

## ai tools

- [x] openai client:
  - https://github.com/sashabaranov/go-openai
  - https://github.com/openai/openai-go 限制的go版本太高
- [x] 通过ai, 翻译内容
- [ ] 支持调用外部命令
- [ ] mcp 支持 https://github.com/mark3labs/mcp-go

## git tools

- [ ] `git_service` 根据当前repo remote 匹配 git service, 并设置相关配置。 比如 github, gitlab, gitea
- [ ] https://github.com/github/github-mcp-server github mcp 工具使用

## fs tools

- [ ] find and clean *.log and more files

## sys tools

- [ ] system info

## 拓展功能

### other tools

- [x] quick jump dir
- [ ] quick json,yml data query
- [ ] 终端自动提示 (auto complete) https://github.com/chzyer/readline

### 扩展插件实现

- https://github.com/d5/tengo Go 的快速脚本语言
- https://github.com/dop251/goja 纯 Go 中的 ECMAScript/JavaScript 引擎
- https://github.com/traefik/yaegi go 语言解释器
- https://github.com/tetratelabs/wazero 纯 Go 的 Wasm 运行时，< 500 KB 增量体积
  - 需要工具先将代码文件转换为 wasm 文件(eg: tinygo, bun)

### 环境工具管理

name: pkg, pkgm, pkgx

- [ ] 配置文件管理
  - 可以迁移，上传配置
  - 一键安装初始化
- [ ] 通过各个平台特有的命令管理软件管理
  - windows: scoop, chocolatey, winget
  - macos: brew, homebrew
  - linux: apt, yum, snap 等

### dev 环境管理

name: `dev env`

- [ ] dev 配置文件管理
- [ ] 安装各种 sdk 例如：php, go, java
  - 配置 sdk 下载地址和安装路径
- [ ] 配置各种环境变量
- [ ] 切换环境，sdk版本等(需要依赖 `kite app shell` 支持)
