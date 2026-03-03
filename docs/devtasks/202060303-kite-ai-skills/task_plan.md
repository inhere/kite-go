# 任务计划：kite ai skills 命令组实现

## 目标
参考 Claude Code 的 skills 管理功能，为 kite-go 实现 `kite ai skills` 命令组，用于统一管理本地的 skills。

---

## 背景研究

### Claude Code Skills 架构
- **目录结构**：
  - 个人：`~/.claude/skills/<skill-name>/SKILL.md`
  - 项目：`.claude/skills/<skill-name>/SKILL.md`
- **SKILL.md 格式**：
  - YAML frontmatter（name, description, disable-model-invocation 等）
  - Markdown 内容（指令）

### kite-go CLI 架构
- **入口**：`cmd/kite/main.go` -> `boot.MustRun(app.App())`
- **框架**：`github.com/gookit/gcli/v3`
- **命令注册**：`internal/cli/boot.go` 的 `addCommands()`
- **现有 AI 命令**：`internal/cli/aicmd/aicmd.go`
- **参考实现**：`internal/cli/xenvcmd/subcmd/tools_cmd.go`

---

## 实现阶段

### Phase 1: 创建数据模型和 Skill 管理服务 [pending]
**目标**：定义 Skill 结构体和管理逻辑

**文件**：
- `internal/cli/aicmd/skills/skill.go` - Skill 结构体定义
- `internal/cli/aicmd/skills/manager.go` - Skill 管理器

**Skill 结构体字段**：
```go
type Skill struct {
    Name        string   `json:"name" yaml:"name"`
    Description string   `json:"description" yaml:"description"`
    Path        string   `json:"path" yaml:"-"`      // SKILL.md 文件路径
    Scope       string   `json:"scope" yaml:"-"`     // user/project
    Content     string   `json:"content" yaml:"-"`   // 原始内容
}
```

**Manager 功能**：
- `ListSkills(scope string) ([]*Skill, error)` - 列出所有 skills
- `GetSkill(name string) (*Skill, error)` - 获取单个 skill
- `CreateSkill(name, description string) error` - 创建新 skill
- `DeleteSkill(name string) error` - 删除 skill
- `EditSkill(name string) error` - 编辑 skill（打开编辑器）
- `ScanSkillsDirs() []string` - 扫描 skills 目录

---

### Phase 2: 实现 skills 子命令 [pending]
**目标**：实现完整命令组

**命令结构**：
```
kite ai skills
├── list (ls)     - 列出所有 skills
├── show          - 显示 skill 详情
├── create (new)  - 创建新 skill
├── edit          - 编辑 skill
├── delete (rm)   - 删除 skill
├── path          - 显示 skills 目录路径
└── open          - 在文件管理器中打开 skills 目录
```

**文件**：
- `internal/cli/aicmd/skills_cmd.go` - 主命令和所有子命令

---

### Phase 3: 集成到 AI 命令组 [pending]
**目标**：将 skills 命令添加到 ai 命令组

**修改文件**：
- `internal/cli/aicmd/aicmd.go` - 添加 SkillsCmd 到 Subs

---

### Phase 4: 测试和文档 [pending]
**目标**：确保功能正常并更新文档

**任务**：
- 编写单元测试
- 更新 README 文档
- 测试命令功能

---

## 预期命令用法

```bash
# 列出所有 skills
kite ai skills list
kite ai skills ls

# 显示 skill 详情
kite ai skills show <name>

# 创建新 skill
kite ai skills create <name> --desc "描述"

# 编辑 skill
kite ai skills edit <name>

# 删除 skill
kite ai skills delete <name>

# 显示 skills 目录路径
kite ai skills path [--scope user|project]

# 在文件管理器中打开
kite ai skills open [--scope user|project]
```

---

## 错误记录

| 错误 | 尝试次数 | 解决方案 |
|------|----------|----------|
| (暂无) | - | - |

---

## 进度追踪

| 阶段 | 状态 | 完成时间 |
|------|------|----------|
| Phase 1 | pending | - |
| Phase 2 | pending | - |
| Phase 3 | pending | - |
| Phase 4 | pending | - |