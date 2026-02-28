# 系统清理功能实现

> 任务日期: 2025-02-13
> 状态: 已完成

## 需求背景

为 `kite-go` 项目实现一个功能完善的系统清理工具，支持清理系统中的临时文件、缓存和无用应用数据，帮助用户释放磁盘空间。

## 功能要点

### 核心功能

1. **并发扫描文件和目录**
   - 使用 goroutine + channel 信号量模式
   - 可配置并发数（默认 5）
   - 支持最大深度限制

2. **清理规则系统**
   - 按类别分组：缓存、日志、临时文件、构建产物、依赖、IDE、系统
   - 风险等级：1（低）、2（中）、3（高）
   - 支持启用/禁用规则
   - 支持多种匹配方式：名称、扩展名、路径模式、年龄、大小

3. **跨平台支持**
   - 通用规则（node_modules, 日志文件, 临时文件等）
   - Windows 特定规则（临时文件夹、更新缓存、缩略图缓存等）
   - macOS 特定规则（系统缓存、Xcode、Homebrew 等）
   - Linux 特定规则（APT/DNF 缓存、系统日志等）

4. **缓存机制**
   - 扫描结果缓存 3 分钟有效
   - 避免短时间内重复扫描
   - 支持手动清除缓存

5. **回收站支持**
   - Windows: 移动到用户目录下的模拟回收站
   - macOS: 移动到 ~/.Trash
   - Linux: 遵循 FreeDesktop.org Trash 规范

6. **报告生成**
   - 支持 Markdown 格式输出
   - 包含扫描统计、分类统计、规则详情
   - 显示警告和错误信息

7. **用户确认**
   - 高风险操作需要确认
   - 支持跳过确认（-y 参数）
   - 支持强制执行（--force 参数）

## 实现方案

### 文件结构

```
internal/service/
├── sys_clean_service.go   # 核心数据类型和服务接口
├── sys_clean_rules.go     # 预设清理规则
├── sys_clean_cache.go     # 缓存管理器
├── sys_clean_scanner.go   # 并发扫描器
├── sys_clean_trash.go     # 回收站管理器
└── sys_clean_report.go    # 报告生成器

internal/cli/syscmd/
└── clean_cmd.go           # 命令行入口
```

### 核心数据结构

```go
// 清理规则
type CleanRule struct {
    Name        string       // 规则名称
    Description string       // 规则描述
    Category    RuleCategory // 类别
    TargetType  TargetType   // 目标类型（文件/目录）
    BasePaths   []string     // 基础路径
    Patterns    []string     // 匹配模式
    FileExts    []string     // 文件扩展名
    NameMatches []string     // 名称匹配
    MaxDepth    int          // 最大深度
    ExcludeDirs []string     // 排除目录
    MinSize     int64        // 最小大小
    MaxAge      int          // 最大年龄（天）
    Platforms   []Platform   // 适用平台
    RiskLevel   int          // 风险等级 1-3
    Enabled     bool         // 是否启用
}

// 扫描结果
type ScanResult struct {
    ID           string         // 缓存ID
    CreatedAt    time.Time      // 创建时间
    ExpiresAt    time.Time      // 过期时间
    Platform     Platform       // 平台
    TotalFiles   int            // 文件总数
    TotalDirs    int            // 目录总数
    TotalTargets int            // 目标总数
    TotalSize    int64          // 总大小
    ScanDuration time.Duration  // 扫描耗时
    Groups       []*TargetGroup // 分组结果
    Errors       []ScanError    // 错误列表
}

// 清理报告
type CleanReport struct {
    GeneratedAt time.Time                       // 生成时间
    Platform    Platform                        // 平台
    Hostname    string                          // 主机名
    Username    string                          // 用户名
    Mode        string                          // 模式
    ScanStats   *ScanStats                      // 扫描统计
    CleanStats  *CleanStats                     // 清理统计
    ByCategory  map[RuleCategory]*CategoryStats // 按类别统计
    ByRule      map[string]*RuleStats           // 按规则统计
    Warnings    []string                        // 警告
    Errors      []string                        // 错误
}
```

### 服务接口

```go
type SysCleanService interface {
    // 配置
    LoadConfig(path string) error
    SetConfig(cfg *CleanConfig)
    GetConfig() *CleanConfig

    // 规则
    GetPresetRules() []*CleanRule
    GetRule(name string) (*CleanRule, bool)
    ListRules() string

    // 扫描
    Scan(ctx context.Context) (*ScanResult, error)
    ScanWithRules(ctx context.Context, ruleNames []string) (*ScanResult, error)

    // 清理
    Preview(ctx context.Context) (*CleanReport, error)
    Clean(ctx context.Context) (*CleanReport, error)

    // 缓存
    LoadCache() (*ScanResult, error)
    SaveCache(result *ScanResult) error
    ClearCache() error
    GetCacheInfo() *CacheInfo

    // 报告
    GenerateReport(result *ScanResult, isExecuted bool) *CleanReport
    ExportMarkdown(report *CleanReport) string
    SaveReport(report *CleanReport, filePath string) error
}
```

## 命令行选项

| 选项 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--config` | `-c` | 配置文件路径 | - |
| `--scan` | `-s` | 扫描目录 | 用户主目录 |
| `--exclude` | `-E` | 排除目录 | .git, .svn, node_modules |
| `--depth` | `-D` | 最大扫描深度 | -1（无限制） |
| `--concurrency` | `-C` | 并发扫描数 | 5 |
| `--rule` | `-R` | 指定规则名称 | - |
| `--category` | | 按类别过滤 | - |
| `--ext` | `-e` | 按文件扩展名过滤 | - |
| `--pattern` | `-p` | 自定义匹配模式 | - |
| `--dry-run` | | 预览模式，不删除文件 | false |
| `--trash` | | 移动到回收站 | false |
| `--yes` | `-y` | 跳过确认 | false |
| `--force` | `-f` | 强制执行高风险操作 | false |
| `--use-cache` | | 使用缓存 | true |
| `--clear-cache` | | 清除缓存 | false |
| `--output` | `-o` | 输出格式 | markdown |
| `--output-file` | `-O` | 输出文件路径 | - |
| `--verbose` | `-v` | 详细输出 | false |
| `--list-rules` | | 列出所有可用规则 | false |

## 预设规则列表

### 通用规则（所有平台）

| 规则名称 | 类别 | 风险 | 说明 |
|----------|------|------|------|
| node_modules | dependency | 1 | Node.js 依赖目录 |
| npm_cache | cache | 1 | npm 缓存目录 |
| yarn_cache | cache | 1 | Yarn 缓存目录 |
| log_files | log | 1 | 日志文件（.log） |
| temp_files | temp | 1 | 临时文件（.tmp, .temp, .bak） |
| editor_backup | temp | 1 | 编辑器备份文件 |
| go_build_cache | build | 1 | Go 构建缓存 |
| python_cache | cache | 1 | Python __pycache__ |
| vscode_cache | ide | 2 | VS Code 缓存 |
| jetbrains_cache | ide | 2 | JetBrains IDE 缓存 |
| dist_build | build | 2 | 构建输出目录 |

### Windows 特定规则

| 规则名称 | 类别 | 风险 | 说明 |
|----------|------|------|------|
| windows_temp | temp | 1 | Windows 临时文件夹 |
| windows_update_cache | cache | 2 | Windows 更新缓存 |
| windows_thumbnail_cache | cache | 1 | 缩略图缓存 |
| browser_cache_windows | cache | 1 | 浏览器缓存 |
| windows_old | system | 3 | Windows.old 文件夹 |

### macOS 特定规则

| 规则名称 | 类别 | 风险 | 说明 |
|----------|------|------|------|
| macos_cache | cache | 1 | macOS 系统缓存 |
| xcode_derived_data | build | 1 | Xcode DerivedData |
| homebrew_cache | cache | 1 | Homebrew 缓存 |
| dmg_residual | temp | 1 | DMG 安装包残留 |
| macos_trash | system | 3 | macOS 废纸篓 |

### Linux 特定规则

| 规则名称 | 类别 | 风险 | 说明 |
|----------|------|------|------|
| linux_temp | temp | 1 | Linux 临时文件 |
| apt_cache | cache | 1 | APT 包管理器缓存 |
| user_cache_linux | cache | 1 | 用户缓存目录 |
| system_logs_linux | log | 2 | 系统日志 |
| linux_trash | system | 3 | Linux 回收站 |

## 配置文件示例

文件包含以下配置项：

┌────────────┬─────────────────────────────────────────────────┐
│  配置分组  │                      说明                       │
├────────────┼─────────────────────────────────────────────────┤
│ 扫描配置   │ scan_dirs, exclude_dirs, max_depth, concurrency │
├────────────┼─────────────────────────────────────────────────┤
│ 规则配置   │ rule_names, categories                          │
├────────────┼─────────────────────────────────────────────────┤
│ 行为配置   │ dry_run, use_trash, yes, force                  │
├────────────┼─────────────────────────────────────────────────┤
│ 缓存配置   │ use_cache, cache_ttl, cache_file                │
├────────────┼─────────────────────────────────────────────────┤
│ 输出配置   │ output_format, output_file, verbose             │
├────────────┼─────────────────────────────────────────────────┤
│ 自定义规则 │ 8 个示例规则，展示各种配置选项                  │
└────────────┴─────────────────────────────────────────────────┘

使用方式：

```bash
kite sys clean --config config/module/sys-clean.yml
```

## 使用示例

```bash
# 列出所有可用规则
kite sys clean --list-rules

# 预览模式（不删除任何文件）
kite sys clean --dry-run

# 预览当前目录
kite sys clean --dry-run -s .

# 执行清理（使用回收站，跳过确认）
kite sys clean --use-trash -y

# 只清理 node_modules
kite sys clean -R node_modules --use-trash -y

# 清理多个规则
kite sys clean -R node_modules -R log_files -R temp_files

# 按类别清理
kite sys clean --category cache --category temp

# 输出到文件
kite sys clean --dry-run -O report.md

# 清除缓存
kite sys clean --clear-cache

# 详细输出
kite sys clean --dry-run -v
```

## 实现结果

### 已完成

- [x] 定义核心数据类型（Platform, TargetType, RuleCategory, CleanRule, CleanTarget 等）
- [x] 实现预设清理规则（通用 + Windows/macOS/Linux 平台规则）
- [x] 实现缓存管理器（3 分钟 TTL，支持加载/保存/清除）
- [x] 实现并发扫描器（goroutine + semaphore 模式）
- [x] 实现回收站管理器（跨平台支持）
- [x] 实现报告生成器（Markdown 格式）
- [x] 组装 SysCleanService 服务层
- [x] 实现命令行入口（完整选项支持）

### 待优化

- [x] JSON 配置文件解析
- [x] YAML 配置文件解析
- [x] JSON 格式报告输出
- [ ] 更精确的 Windows 回收站 API 调用（SHFileOperation）
- [ ] 单元测试覆盖

## 技术说明

### 并发扫描策略

使用 `sync.WaitGroup` + channel 信号量模式：
1. 每个扫描目录启动一个 goroutine
2. 通过 semaphore channel 控制最大并发数
3. 使用 targetChan 收集扫描结果
4. 使用 errorChan 收集错误
5. 结果使用 mutex 保护

### 缓存机制

- 缓存文件路径: `~/.kite-go/tmp/sys-clean/scan-cache.json`
- 默认 TTL: 3 分钟
- 包含完整扫描结果和元信息

### 回收站实现

- **Windows**: 移动到 `~/.kite-go/trash/` 目录（简化实现）
- **macOS**: 移动到 `~/.Trash/` 系统废纸篓
- **Linux**: 使用 FreeDesktop.org Trash 规范，移动到 `~/.local/share/Trash/files/`

## 参考资料

- [FreeDesktop.org Trash Specification](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html)
- [filepath.Match - Go Documentation](https://pkg.go.dev/path/filepath#Match)
