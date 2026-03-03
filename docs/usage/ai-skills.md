# AI Skills 管理

`kite ai skills` 命令用于统一管理本地的 AI Skills，兼容 Claude Code 的 Skills 规范。

## 目录结构

Skills 存储在两个级别：

| 级别 | 路径 | 说明 |
|------|------|------|
| User | `~/.claude/skills/<skill-name>/SKILL.md` | 用户级别，跨项目共享 |
| Project | `./.claude/skills/<skill-name>/SKILL.md` | 项目级别，仅当前项目可用 |

每个 Skill 是一个目录，包含 `SKILL.md` 文件作为入口点。

## 如述

```bash
# 查看所有子命令
kite ai skills --help

# 列出所有 skills
kite ai skills list

# 创建新 skill
kite ai skills create my-skill

# 显示 skill 详情
kite ai skills show my-skill

# 删除 skill
kite ai skills delete my-skill
```

## 子命令

### list - 列出所有 Skills

列出所有可用的 Skills，支持按 scope 过滤。

```bash
# 列出所有 skills
kite ai skills list

# 只列出用户级别 skills
kite ai skills list --scope user

# 只列出项目级别 skills
kite ai skills list --scope project
```

**选项**：

- `--scope, -s` - 过滤范围： `user`、`project` 或 `all`（默认）

**别名**: `ls`

### show - 显示 Skill 详情

显示指定 Skill 的详细信息，包括 frontmatter 和内容。

```bash
kite ai skills show <name>
```

**示例**：

```bash
kite ai skills show my-skill
```

### create - 创建新 Skill

创建一个新的 Skill，会自动生成标准模板文件。

```bash
kite ai skills create <name> [options]
```

**选项**：

- `--desc, -d` - Skill 描述
- `--scope, -s` - 创建位置： `user`（默认）或 `project`

**别名**: `new`, `add`

**示例**：

```bash
# 创建用户级别 skill
kite ai skills create my-skill --desc "My custom skill"

# 创建项目级别 skill
kite ai skills create project-skill --scope project
```

### edit - 编辑 Skill

在默认编辑器中打开 Skill 文件进行编辑。

```bash
kite ai skills edit <name>
```

编辑器选择顺序：
1. 环境变量 `EDITOR`
2. 环境变量 `VISUAL`
3. Windows: `notepad`，其他系统: `vim`

**示例**：

```bash
kite ai skills edit my-skill
```

### delete - 删除 Skill

删除指定的 Skill。

```bash
kite ai skills delete <name> [options]
```

**选项**：

- `--force, -f` - 跳过确认直接删除

**别名**: `rm`, `remove`

> **注意**: `-f` 选项必须在参数之前。

**示例**：

```bash
# 删除（需要确认）
kite ai skills delete my-skill

# 强制删除（无需确认）
kite ai skills delete -f my-skill
```

### path - 显示 Skills 目录路径

显示 Skills 存储目录的路径。

```bash
# 显示所有路径
kite ai skills path

# 只显示用户级别路径
kite ai skills path --scope user

# 只显示项目级别路径
kite ai skills path --scope project
```

**选项**：

- `--scope, -s` - 显示指定范围的路径

### open - 在文件管理器中打开

在系统文件管理器中打开 Skills 目录。

```bash
# 打开用户级别目录
kite ai skills open

# 打开项目级别目录
kite ai skills open --scope project
```

**选项**：

- `--scope, -s` - 打开指定范围的目录

## SKILL.md 文件格式

每个 Skill 的 `SKILL.md` 文件包含 YAML frontmatter 和 Markdown 内容。

```yaml
---
name: skill-name
description: What this skill does and when to use it
argument-hint: [arguments]
disable-model-invocation: false
user-invocable: true
allowed-tools:
  - Read
  - Grep
model: claude-sonnet-4-20250514
---

# Skill Title

Instructions for the skill go here.

## Usage

Describe how to use this skill.

## Examples

Provide examples if helpful.
```

### Frontmatter 字段说明

| 字段 | 必需 | 说明 |
|------|------|------|
| `name` | 否 | Skill 名称，默认使用目录名 |
| `description` | 推荐 | 描述，用于自动加载判断 |
| `argument-hint` | 否 | 参数提示，如 `[issue-number]` |
| `disable-model-invocation` | 否 | 设为 `true` 禁止 AI 自动加载 |
| `user-invocable` | 否 | 设为 `false` 从菜单隐藏 |
| `allowed-tools` | 否 | 允许使用的工具列表 |
| `model` | 否 | 指定使用的模型 |

## 使用场景

### 场景 1：创建项目特定的 Skill

```bash
# 在项目中创建 skill
kite ai skills create code-review --scope project --desc "Code review guidelines"

# 编辑 skill
kite ai skills edit code-review
```

### 场景 2：管理共享 Skills

```bash
# 创建用户级别 skill
kite ai skills create git-conventions --desc "Git commit conventions"

# 列出所有 skills
kite ai skills list

# 删除不再需要的 skill
kite ai skills delete -f old-skill
```

### 场景 3：快速查看 Skills 位置

```bash
# 查看路径
kite ai skills path

# 在文件管理器中打开
kite ai skills open
```

## 与 Claude Code 兼容

`kite ai skills` 命令完全兼容 Claude Code 的 Skills 规范：
- 使用相同的目录结构
- 支持相同的 SKILL.md 格式
- 可以共享 Skills

## 命令速查表

| 命令 | 说明 |
|------|------|
| `kite ai skills list` | 列出所有 skills |
| `kite ai skills show <name>` | 显示 skill 详情 |
| `kite ai skills create <name>` | 创建新 skill |
| `kite ai skills edit <name>` | 编辑 skill |
| `kite ai skills delete <name>` | 删除 skill |
| `kite ai skills path` | 显示目录路径 |
| `kite ai skills open` | 打开目录 |
