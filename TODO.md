# TODO

- [ ] go plugin 支持 https://github.com/hashicorp/go-plugin
- [ ] 表达式支持
    - https://github.com/expr-lang/expr Go 的表达式语言和表达式评估
    - https://github.com/hashicorp/go-bexpr Go 中的通用布尔表达式评估
    - https://github.com/ganigeorgiev/fexpr 支持解析类似 SQL 表达式 eg: `id=123 && status='active'`

## kite app

- [ ] generate metadata file for kite

## run anything

- [x] scripts manage and run 

## backend serve

- [ ] support backend server: `kited` `kite app:server`
- [ ] cache metadata, provide quick command search

## job/task serve

- [ ] start a background server, listen an port/sock
- [ ] can delivery task by tcp connection

## http tools

- [ ] provide web UI for some operation
- [ ] can delivery task by web page
- [ ] http benchmark tool
- [ ] send http request tool. like curl, ide-http-client

## other tools

- [x] quick jump dir
- [ ] quick json,yml data query
- [ ] 终端自动提示 (auto complete) https://github.com/chzyer/readline

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

## extra command

支持调用外部命令

- [ ] 调用 php 实现的命令
