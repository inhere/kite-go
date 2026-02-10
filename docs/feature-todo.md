# Todo

使用 emoji 表示任务状态

- 进行中emoji: 🚧 / 🔄 
- 已实现emoji: ✅
- 待实现emoji: ⏳

## kite ai claude 管理命令1 ✅

需求：在 `internal/cli/aicmd` 新增一个命令组 `claude_cmd`, 用于快速管理 claude code 的一些设置和信息。

`kite ai claude -h` 提供命令：

- api   配置当前要使用的 claude api server 信息等

### 补充信息

- claude 命令有自己的配置文件 ai-claude.yaml 运行时自动加载，provider map 配置在这里面

### 子命令: api

```bash
# 设置要使用的server api
$ kite ai claude api --use glm|minimax|kimi|claude --shell pwsh|bash
```

- `--use`: 指定使用哪个模型，设置后将会修改 `~/.claude/config.json` 文件
- `--shell`: 指定当前使用的 shell，指定后将会输出当前shell的环境变量，用于在 shell 中使用
- `--write, -w` 是否写入 claude 配置文件，默认不写入

### claude 配置文件参考

```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "{api-auth-token}",
    "ANTHROPIC_BASE_URL": "{api-server-address}"
  },
  "includeCoAuthoredBy": false
}
```

## kite ai skills 管理命令

提供命令：

```bash
kite ai skills list
kite ai skills list --keyword xxx
```

