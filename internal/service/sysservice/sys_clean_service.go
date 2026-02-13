package sysservice

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/inhere/kite-go/internal/service/sysservice/sysclean"
)

// SysCleanService 系统清理服务接口
type SysCleanService interface {
	// 配置
	LoadConfig(path string) error
	SetConfig(cfg *sysclean.CleanConfig)
	GetConfig() *sysclean.CleanConfig

	// 规则
	GetPresetRules() []*sysclean.CleanRule
	GetRule(name string) (*sysclean.CleanRule, bool)
	ListRules() string

	// 扫描
	Scan(ctx context.Context) (*sysclean.ScanResult, error)
	ScanWithRules(ctx context.Context, ruleNames []string) (*sysclean.ScanResult, error)

	// 清理
	Preview(ctx context.Context) (*sysclean.CleanReport, error)
	Clean(ctx context.Context) (*sysclean.CleanReport, error)

	// 缓存
	LoadCache() (*sysclean.ScanResult, error)
	SaveCache(result *sysclean.ScanResult) error
	ClearCache() error
	GetCacheInfo() *CacheInfo

	// 报告
	GenerateReport(result *sysclean.ScanResult, isExecuted bool) *sysclean.CleanReport
	ExportMarkdown(report *sysclean.CleanReport) string
	SaveReport(report *sysclean.CleanReport, filePath string) error
}

// sysCleanService 系统清理服务实现
type sysCleanService struct {
	config       *sysclean.CleanConfig
	presetRules  []*sysclean.CleanRule
	cacheManager *sysclean.CacheManager
	trashManager sysclean.TrashManager
}

// NewSysCleanService 创建系统清理服务
func NewSysCleanService(config *sysclean.CleanConfig) SysCleanService {
	if config == nil {
		config = sysclean.DefaultCleanConfig()
	}

	// 初始化缓存管理器
	cacheFile := config.CacheFile
	if cacheFile == "" {
		homeDir, _ := os.UserHomeDir()
		cacheFile = homeDir + "/.kite-go/tmp/sys-clean/scan-cache.json"
	}

	return &sysCleanService{
		config:       config,
		presetRules:  sysclean.GetPresetRules(),
		cacheManager: sysclean.NewCacheManager(cacheFile, config.CacheTTL),
		trashManager: sysclean.NewTrashManager(),
	}
}

// SysClean 获取系统清理服务
func SysClean() SysCleanService {
	return NewSysCleanService(nil)
}

// SysCleanWithConfig 使用配置获取系统清理服务
func SysCleanWithConfig(cfg *sysclean.CleanConfig) SysCleanService {
	return NewSysCleanService(cfg)
}

// Scan 实现服务接口的扫描
func (s *sysCleanService) Scan(ctx context.Context) (*sysclean.ScanResult, error) {
	// 检查是否使用缓存
	if s.config.UseCache {
		cached, err := s.cacheManager.Load()
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	// 获取启用的规则
	rules := sysclean.FilterEnabledRules(s.presetRules)
	rules = sysclean.FilterRulesByPlatform(rules, sysclean.CurrentPlatform())

	// 应用规则名称过滤
	if len(s.config.RuleNames) > 0 {
		rules = sysclean.FilterRulesByNames(rules, s.config.RuleNames)
	}

	// 应用类别过滤
	if len(s.config.Categories) > 0 {
		rules = sysclean.FilterRulesByCategories(rules, s.config.Categories)
	}

	// 获取扫描目录
	scanDirs := s.config.ScanDirs
	if len(scanDirs) == 0 {
		homeDir, _ := os.UserHomeDir()
		scanDirs = []string{homeDir}
	}

	// 创建扫描器并执行扫描
	scanner := sysclean.NewScanner(s.config, rules)
	result, err := scanner.Scan(ctx, scanDirs)
	if err != nil {
		return nil, err
	}

	// 保存配置快照
	result.Config = s.config

	// 缓存结果
	if s.config.UseCache {
		_ = s.cacheManager.Save(result)
	}

	return result, nil
}

// ScanWithRules 使用指定规则扫描
func (s *sysCleanService) ScanWithRules(ctx context.Context, ruleNames []string) (*sysclean.ScanResult, error) {
	// 临时保存原规则名称
	originalNames := s.config.RuleNames
	s.config.RuleNames = ruleNames
	defer func() {
		s.config.RuleNames = originalNames
	}()

	return s.Scan(ctx)
}

// LoadConfig 加载配置文件
func (s *sysCleanService) LoadConfig(cfgFile string) error {
	cfg := config.NewGeneric("ai-config", config.WithTagName("json"))
	cfg.AddDriver(yaml.Driver)

	if err := cfg.LoadFiles(cfgFile); err != nil {
		return fmt.Errorf("failed to load config file %s: %w", cfgFile, err)
	}
	if err := cfg.Decode(&s.config); err != nil {
		return fmt.Errorf("failed to decode AI config: %w", err)
	}
	cfg.ClearAll()

	// 添加自定义规则到预设规则
	if len(s.config.Rules) > 0 {
		s.presetRules = append(s.presetRules, s.config.Rules...)
	}
	return nil
}

// GetConfig 获取配置
func (s *sysCleanService) GetConfig() *sysclean.CleanConfig {
	return s.config
}

// SetConfig 设置配置
func (s *sysCleanService) SetConfig(cfg *sysclean.CleanConfig) {
	s.config = cfg
}

// GetPresetRules 获取预设规则
func (s *sysCleanService) GetPresetRules() []*sysclean.CleanRule {
	return s.presetRules
}

// GetRule 获取指定规则
func (s *sysCleanService) GetRule(name string) (*sysclean.CleanRule, bool) {
	for _, rule := range s.presetRules {
		if rule.Name == name {
			return rule, true
		}
	}
	return nil, false
}

// LoadCache 实现服务接口的缓存加载
func (s *sysCleanService) LoadCache() (*sysclean.ScanResult, error) {
	return s.cacheManager.Load()
}

// SaveCache 实现服务接口的缓存保存
func (s *sysCleanService) SaveCache(result *sysclean.ScanResult) error {
	return s.cacheManager.Save(result)
}

// ClearCache 实现服务接口的缓存清除
func (s *sysCleanService) ClearCache() error {
	return s.cacheManager.Clear()
}

// CacheInfo 缓存信息
type CacheInfo struct {
	Exists       bool          `json:"exists"`
	FilePath     string        `json:"file_path"`
	TTL          time.Duration `json:"ttl"`
	Remaining    time.Duration `json:"remaining"`
	IsExpired    bool          `json:"is_expired"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
	ExpiresAt    time.Time     `json:"expires_at,omitempty"`
	TotalTargets int           `json:"total_targets,omitempty"`
	TotalSize    int64         `json:"total_size,omitempty"`
}

// GetCacheInfo 获取缓存信息
func (s *sysCleanService) GetCacheInfo() *CacheInfo {
	info := &CacheInfo{
		Exists:    s.cacheManager.Exists(),
		FilePath:  s.cacheManager.GetFilePath(),
		TTL:       s.cacheManager.GetTTL(),
		IsExpired: s.cacheManager.IsExpired(),
		Remaining: s.cacheManager.RemainingTime(),
	}

	if info.Exists && !info.IsExpired {
		result, err := s.cacheManager.Load()
		if err == nil {
			info.CreatedAt = result.CreatedAt
			info.ExpiresAt = result.ExpiresAt
			info.TotalTargets = result.TotalTargets
			info.TotalSize = result.TotalSize
		}
	}

	return info
}
