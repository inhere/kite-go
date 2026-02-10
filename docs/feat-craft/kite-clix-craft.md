# Kite Clix 命令功能草稿

为 kite cli新增实现一个命令 `xcli` 用于执行项目目录下的 `kite.xcli.yaml` 文件中定义的命令。
用于简单快速的为一个应用项目实现一个命令行工具。

## 命令定义

命令定义参考：

```yaml
# 内置设置
__settings:
  author: "kite-clix"
  version: "1.0.0"
  description: "当前命令应用的描述"
  env:
    KEY1: "value1"

demo-cmd:
  description: "demo cmd 命令的描述"
  # run 配置运行一段 bash 脚本
  run: |
    echo "from demo-cmd"

demo-cmd2:
  description: "demo cmd2 命令的描述"
  # script 配置运行一个脚本文件
  script: ./build-script.sh
```

说明：

- `__settings` 配置项目信息和环境变量等
  - version, description 用于显示在 `kite xcli -l` 命令中。
- `<command>` 配置一个命令，包含命令的描述和运行方式。
  - `run` 配置运行一段 bash 脚本。与 `script` 配置互斥。
  - `script` 配置运行一个脚本文件。

## 命令使用

- `kite xcli` 或者 `kite xcli -h|--help` 显示 `xcli` 自己的帮助信息，列出所有内部命令。
- `kite xcli -l|--list` 会自动检测当前目录下是否存在 `kite.xcli.yaml` 文件，如果存在则读取并列出命令。
- `kite xcli -c|--create` 创建一个命令定义文件 `kite.xcli.yaml`，并自动填充示例命令定义。
- `kite xcli <command>` 检测并运行配置文件定义的 `<command>` 命令。
