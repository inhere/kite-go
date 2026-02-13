package sysclean

import (
	"os"
	"runtime"
	"time"
)

// Platform 当前运行平台
type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformDarwin  Platform = "darwin"
	PlatformLinux   Platform = "linux"
	PlatformAll     Platform = "all"
)

// CurrentPlatform 获取当前平台
func CurrentPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return PlatformWindows
	case "darwin":
		return PlatformDarwin
	case "linux":
		return PlatformLinux
	default:
		return PlatformLinux
	}
}

// TargetType 目标类型
type TargetType string

const (
	TargetTypeFile    TargetType = "file"
	TargetTypeDir     TargetType = "dir"
	TargetTypeBoth    TargetType = "both"
	TargetTypePattern TargetType = "pattern"
)

// RuleCategory 规则类别
type RuleCategory string

const (
	CategoryCache      RuleCategory = "cache"
	CategoryLog        RuleCategory = "log"
	CategoryTemp       RuleCategory = "temp"
	CategoryBuild      RuleCategory = "build"
	CategoryDependency RuleCategory = "dependency"
	CategoryIDE        RuleCategory = "ide"
	CategorySystem     RuleCategory = "system"
	CategoryCustom     RuleCategory = "custom"
)

// DisplayName 获取类别显示名称
func (rc RuleCategory) DisplayName() string {
	switch rc {
	case CategoryCache:
		return "缓存文件"
	case CategoryLog:
		return "日志文件"
	case CategoryTemp:
		return "临时文件"
	case CategoryBuild:
		return "构建产物"
	case CategoryDependency:
		return "依赖目录"
	case CategoryIDE:
		return "IDE 缓存"
	case CategorySystem:
		return "系统文件"
	case CategoryCustom:
		return "自定义"
	default:
		return string(rc)
	}
}


// CleanRule 清理规则
type CleanRule struct {
	// 基本属性
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Category    RuleCategory `json:"category" yaml:"category"`
	TargetType  TargetType `json:"target_type" yaml:"target_type"`

	// 匹配规则
	BasePaths   []string `json:"base_paths,omitempty" yaml:"base_paths,omitempty"`
	Patterns    []string `json:"patterns,omitempty" yaml:"patterns,omitempty"`
	FileExts    []string `json:"file_exts,omitempty" yaml:"file_exts,omitempty"`
	NameMatches []string `json:"name_matches,omitempty" yaml:"name_matches,omitempty"`

	// 扫描配置
	MaxDepth    int      `json:"max_depth,omitempty" yaml:"max_depth,omitempty"`
	ExcludeDirs []string `json:"exclude_dirs,omitempty" yaml:"exclude_dirs,omitempty"`

	// 条件过滤
	MinSize int64 `json:"min_size,omitempty" yaml:"min_size,omitempty"` // 最小文件大小（字节）
	MaxAge  int   `json:"max_age,omitempty" yaml:"max_age,omitempty"`   // 最大文件年龄（天）

	// 平台和风险
	Platforms []Platform `json:"platforms,omitempty" yaml:"platforms,omitempty"`
	RiskLevel int        `json:"risk_level" yaml:"risk_level"` // 1-低, 2-中, 3-高

	// 清理行为
	ConfirmMsg string `json:"confirm_msg,omitempty" yaml:"confirm_msg,omitempty"`
	UseTrash   bool   `json:"use_trash,omitempty" yaml:"use_trash,omitempty"`
	Recursive  bool   `json:"recursive,omitempty" yaml:"recursive,omitempty"`
	Enabled    bool   `json:"enabled" yaml:"enabled"`
}

// CleanTarget 扫描目标
type CleanTarget struct {
	Path      string       `json:"path"`
	Name      string       `json:"name"`
	Size      int64        `json:"size"`
	IsDir     bool         `json:"is_dir"`
	ModTime   time.Time    `json:"mod_time"`
	RuleName  string       `json:"rule_name"`
	Category  RuleCategory `json:"category"`
	RiskLevel int          `json:"risk_level"`
}

// ScanError 扫描错误
type ScanError struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// TargetGroup 按规则分组的目标
type TargetGroup struct {
	RuleName  string         `json:"rule_name"`
	Category  RuleCategory   `json:"category"`
	RiskLevel int            `json:"risk_level"`
	Targets   []*CleanTarget `json:"targets"`
	TotalSize int64          `json:"total_size"`
	FileCount int            `json:"file_count"`
	DirCount  int            `json:"dir_count"`
}

// ScanStats 扫描统计
type ScanStats struct {
	TotalFiles   int           `json:"total_files"`
	TotalDirs    int           `json:"total_dirs"`
	TotalTargets int           `json:"total_targets"`
	TotalSize    int64         `json:"total_size"`
	Duration     time.Duration `json:"duration"`
}

// CleanStats 清理统计
type CleanStats struct {
	DeletedFiles  int           `json:"deleted_files"`
	DeletedDirs   int           `json:"deleted_dirs"`
	FailedCount   int           `json:"failed_count"`
	FreedSpace    int64         `json:"freed_space"`
	FailedSpace   int64         `json:"failed_space"`
	Duration      time.Duration `json:"duration"`
	UseTrash      bool          `json:"use_trash"`
}

// CategoryStats 分类统计
type CategoryStats struct {
	Category  RuleCategory `json:"category"`
	FileCount int          `json:"file_count"`
	DirCount  int          `json:"dir_count"`
	TotalSize int64        `json:"total_size"`
}

// RuleStats 规则统计
type RuleStats struct {
	RuleName  string `json:"rule_name"`
	FileCount int    `json:"file_count"`
	DirCount  int    `json:"dir_count"`
	TotalSize int64  `json:"total_size"`
}

// ScanResult 扫描结果
type ScanResult struct {
	// 缓存元信息
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`

	// 扫描环境
	Platform Platform `json:"platform"`

	// 配置快照
	Config *CleanConfig `json:"config,omitempty"`

	// 统计
	TotalFiles   int           `json:"total_files"`
	TotalDirs    int           `json:"total_dirs"`
	TotalTargets int           `json:"total_targets"`
	TotalSize    int64         `json:"total_size"`
	ScanDuration time.Duration `json:"scan_duration"`

	// 分组结果
	Groups []*TargetGroup `json:"groups"`

	// 错误列表
	Errors []ScanError `json:"errors,omitempty"`
}

// CleanReport 清理报告
type CleanReport struct {
	GeneratedAt time.Time `json:"generated_at"`
	Platform    Platform  `json:"platform"`
	Hostname    string    `json:"hostname"`
	Username    string    `json:"username"`
	Mode        string    `json:"mode"` // dry-run, preview, clean

	// 扫描统计
	ScanStats *ScanStats `json:"scan_stats,omitempty"`

	// 清理统计
	CleanStats *CleanStats `json:"clean_stats,omitempty"`

	// 分类统计
	ByCategory map[RuleCategory]*CategoryStats `json:"by_category"`

	// 规则统计
	ByRule map[string]*RuleStats `json:"by_rule"`

	// 警告和错误
	Warnings []string `json:"warnings,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

// CleanConfig 清理配置
type CleanConfig struct {
	// 扫描配置
	ScanDirs    []string `json:"scan_dirs,omitempty" yaml:"scan_dirs,omitempty"`
	Concurrency int      `json:"concurrency,omitempty" yaml:"concurrency,omitempty"` // 默认 5
	MaxDepth    int      `json:"max_depth,omitempty" yaml:"max_depth,omitempty"`     // -1 无限制
	ExcludeDirs []string `json:"exclude_dirs,omitempty" yaml:"exclude_dirs,omitempty"`

	// 规则配置
	Rules      []*CleanRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	RuleNames  []string     `json:"rule_names,omitempty" yaml:"rule_names,omitempty"`
	Categories []RuleCategory `json:"categories,omitempty" yaml:"categories,omitempty"`

	// 行为配置
	DryRun   bool     `json:"dry_run,omitempty" yaml:"dry_run,omitempty"`
	UseTrash bool     `json:"use_trash,omitempty" yaml:"use_trash,omitempty"`
	Force    bool     `json:"force,omitempty" yaml:"force,omitempty"`
	Yes      bool     `json:"yes,omitempty" yaml:"yes,omitempty"`
	FileExts []string `json:"file_exts,omitempty" yaml:"file_exts,omitempty"`
	Pattern  string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// 缓存配置
	UseCache bool          `json:"use_cache,omitempty" yaml:"use_cache,omitempty"`
	CacheTTL time.Duration `json:"cache_ttl,omitempty" yaml:"cache_ttl,omitempty"` // 默认 3 分钟
	CacheFile string       `json:"cache_file,omitempty" yaml:"cache_file,omitempty"`

	// 输出配置
	OutputFormat string `json:"output_format,omitempty" yaml:"output_format,omitempty"` // text, markdown, json
	OutputFile   string `json:"output_file,omitempty" yaml:"output_file,omitempty"`
	Verbose      bool   `json:"verbose,omitempty" yaml:"verbose,omitempty"`
}

// DefaultCleanConfig 返回默认配置
func DefaultCleanConfig() *CleanConfig {
	homeDir, _ := os.UserHomeDir()
	return &CleanConfig{
		ScanDirs:    []string{homeDir},
		Concurrency: 5,
		MaxDepth:    -1,
		ExcludeDirs: []string{".git", ".svn", "node_modules"},
		UseCache:    true,
		CacheTTL:    3 * time.Minute,
		DryRun:      true,
		OutputFormat: "markdown",
	}
}
