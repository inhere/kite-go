# 研究发现：Skills 命令实现

## 1. Claude Code Skills 规范

### 目录结构
```
~/.claude/skills/          # 个人 skills（跨项目）
├── my-skill/
│   ├── SKILL.md           # 必需：主指令文件
│   ├── reference.md       # 可选：详细参考文档
│   ├── examples/          # 可选：示例目录
│   └── scripts/           # 可选：脚本文件
│
.claude/skills/            # 项目 skills（仅当前项目）
├── project-skill/
│   └── SKILL.md
```

### SKILL.md 格式
```yaml
---
name: skill-name                    # 可选，默认使用目录名
description: What this skill does   # 推荐，用于自动加载判断
argument-hint: [issue-number]       # 可选，参数提示
disable-model-invocation: true      # 可选，禁止自动加载
user-invocable: false               # 可选，从菜单隐藏
allowed-tools: Read, Grep           # 可选，允许的工具
model: claude-sonnet-4-20250514     # 可选，指定模型
context: fork                       # 可选，在子代理中运行
agent: Explore                      # 可选，子代理类型
---

# Skill 内容
指令内容...
```

### Frontmatter 字段
| 字段 | 必需 | 描述 |
|------|------|------|
| `name` | 否 | 显示名称，默认使用目录名 |
| `description` | 推荐 | 描述，用于自动加载判断 |
| `argument-hint` | 否 | 参数提示 |
| `disable-model-invocation` | 否 | 禁止 Claude 自动加载 |
| `user-invocable` | 否 | 从 `/` 菜单隐藏 |
| `allowed-tools` | 否 | 允许的工具列表 |
| `model` | 否 | 指定使用的模型 |
| `context` | 否 | 设为 `fork` 在子代理中运行 |
| `agent` | 否 | 子代理类型 |

---

## 2. kite-go CLI 架构

### 命令注册流程
```
cmd/kite/main.go
    └── boot.MustRun(app.App())
            └── boot.Boot(ka)
                    └── BootCli()
                            └── cli.Boot(cliApp)
                                    └── addCommands(cli)
                                            └── cli.Add(aicmd.AICommand, ...)
```

### 命令定义模式
```go
var MyCmd = &gcli.Command{
    Name:    "mycmd",
    Desc:    "命令描述",
    Aliases: []string{"mc"},
    Subs: []*gcli.Command{
        SubCmd1,
        SubCmd2,
    },
    Config: func(c *gcli.Command) {
        // 配置选项和参数
    },
    Func: func(c *gcli.Command, args []string) error {
        // 命令逻辑
    },
}
```

### 参考文件
| 文件 | 用途 |
|------|------|
| `internal/cli/boot.go` | 命令注册入口 |
| `internal/cli/aicmd/aicmd.go` | AI 命令组 |
| `internal/cli/xenvcmd/subcmd/tools_cmd.go` | 工具管理命令（install/list/show 等） |
| `internal/cli/xenvcmd/subcmd/config_cmd.go` | 配置管理命令 |

---

## 3. 关键设计决策

### Skills 目录
- 用户级：`~/.claude/skills/`（兼容 Claude Code）
- 项目级：`./.claude/skills/`（项目目录下）

### 需要实现的功能
1. **list** - 列出所有 skills（支持 scope 过滤）
2. **show** - 显示 skill 详情（内容、frontmatter 等）
3. **create** - 创建新 skill（生成模板文件）
4. **edit** - 编辑 skill（打开编辑器）
5. **delete** - 删除 skill
6. **path** - 显示 skills 目录路径
7. **open** - 在文件管理器中打开

### Skill 结构体
```go
type Skill struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Path        string            `json:"path"`
    Scope       string            `json:"scope"`       // user/project
    Content     string            `json:"content"`
    Frontmatter map[string]any    `json:"frontmatter"`
}
```