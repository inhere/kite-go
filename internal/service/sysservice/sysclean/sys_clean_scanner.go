package sysclean

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// Scanner 并发扫描器
type Scanner struct {
	config    *CleanConfig
	rules     []*CleanRule
	semaphore chan struct{}
	wg        sync.WaitGroup
	result    *ScanResult
	resultMu  sync.Mutex
}

// NewScanner 创建扫描器
func NewScanner(config *CleanConfig, rules []*CleanRule) *Scanner {
	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}
	return &Scanner{
		config:    config,
		rules:     rules,
		semaphore: make(chan struct{}, concurrency),
		result: &ScanResult{
			ID:        generateID(),
			CreatedAt: time.Now(),
			Platform:  CurrentPlatform(),
			Groups:    make([]*TargetGroup, 0),
			Errors:    make([]ScanError, 0),
		},
	}
}

// generateID 生成唯一ID
func generateID() string {
	return time.Now().Format("20060102150405")
}

// Scan 执行扫描
func (s *Scanner) Scan(ctx context.Context, dirs []string) (*ScanResult, error) {
	startTime := time.Now()

	// 初始化结果
	s.result = &ScanResult{
		ID:        generateID(),
		CreatedAt: startTime,
		Platform:  CurrentPlatform(),
		Groups:    make([]*TargetGroup, 0),
		Errors:    make([]ScanError, 0),
	}

	// 初始化规则分组
	ruleGroups := make(map[string]*TargetGroup)
	for _, rule := range s.rules {
		ruleGroups[rule.Name] = &TargetGroup{
			RuleName:  rule.Name,
			Category:  rule.Category,
			RiskLevel: rule.RiskLevel,
			Targets:   make([]*CleanTarget, 0),
		}
	}

	// 创建通道
	targetChan := make(chan *CleanTarget, 100)
	errorChan := make(chan ScanError, 50)
	doneChan := make(chan struct{})

	// 启动结果收集器
	go func() {
		for {
			select {
			case target := <-targetChan:
				s.addTarget(target, ruleGroups)
			case err := <-errorChan:
				s.addError(err)
			case <-doneChan:
				return
			}
		}
	}()

	// 并发扫描各个目录
	for _, dir := range dirs {
		s.wg.Add(1)
		go s.scanDir(ctx, dir, targetChan, errorChan)
	}

	// 等待所有扫描完成
	s.wg.Wait()
	close(doneChan)

	// 收集结果
	s.result.ScanDuration = time.Since(startTime)
	for _, group := range ruleGroups {
		if len(group.Targets) > 0 {
			s.result.Groups = append(s.result.Groups, group)
			s.result.TotalTargets += len(group.Targets)
			s.result.TotalSize += group.TotalSize
			s.result.TotalFiles += group.FileCount
			s.result.TotalDirs += group.DirCount
		}
	}

	return s.result, nil
}

// addTarget 添加扫描目标
func (s *Scanner) addTarget(target *CleanTarget, ruleGroups map[string]*TargetGroup) {
	s.resultMu.Lock()
	defer s.resultMu.Unlock()

	group, ok := ruleGroups[target.RuleName]
	if !ok {
		return
	}

	group.Targets = append(group.Targets, target)
	group.TotalSize += target.Size
	if target.IsDir {
		group.DirCount++
	} else {
		group.FileCount++
	}
}

// addError 添加扫描错误
func (s *Scanner) addError(err ScanError) {
	s.resultMu.Lock()
	defer s.resultMu.Unlock()
	s.result.Errors = append(s.result.Errors, err)
}

// scanDir 扫描目录
func (s *Scanner) scanDir(ctx context.Context, dir string, targetChan chan<- *CleanTarget, errorChan chan<- ScanError) {
	defer s.wg.Done()

	// 获取信号量
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return
	default:
	}

	// 展开路径中的 ~ 和环境变量
	dir = fsutil.ExpandPath(dir)

	// 遍历目录
	err := fsutil.FindInDir(dir, func(path string, ent fs.DirEntry) error {
		// 检查上下文
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 获取文件信息
		info, err := ent.Info()
		if err != nil {
			return nil // 忽略无法获取信息的文件
		}

		// 检查排除目录
		if s.isExcluded(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 匹配规则
		for _, rule := range s.rules {
			if s.matchesRule(path, info, rule) {
				targetChan <- &CleanTarget{
					Path:      path,
					Name:      info.Name(),
					Size:      s.getSize(path, info),
					IsDir:     info.IsDir(),
					ModTime:   info.ModTime(),
					RuleName:  rule.Name,
					Category:  rule.Category,
					RiskLevel: rule.RiskLevel,
				}
				break // 每个路径只匹配一个规则
			}
		}

		return nil
	})

	if err != nil && err != filepath.SkipDir {
		errorChan <- ScanError{
			Path:  dir,
			Error: err.Error(),
		}
	}
}

// isExcluded 检查路径是否在排除列表中
func (s *Scanner) isExcluded(path string) bool {
	for _, exclude := range s.config.ExcludeDirs {
		if strings.Contains(path, exclude) {
			return true
		}
	}
	return false
}

// getMaxDepth 获取最大深度
func (s *Scanner) getMaxDepth() int {
	if s.config.MaxDepth < 0 {
		return 0 // 0 表示无限制
	}
	return s.config.MaxDepth
}

// getSize 获取文件/目录大小
func (s *Scanner) getSize(path string, info os.FileInfo) int64 {
	if info.IsDir() {
		// 对于目录，计算总大小
		var size int64
		filepath.Walk(path, func(_ string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				size += fi.Size()
			}
			return nil
		})
		return size
	}
	return info.Size()
}

// matchesRule 检查路径是否匹配规则
func (s *Scanner) matchesRule(path string, info os.FileInfo, rule *CleanRule) bool {
	// 检查目标类型
	switch rule.TargetType {
	case TargetTypeFile:
		if info.IsDir() {
			return false
		}
	case TargetTypeDir:
		if !info.IsDir() {
			return false
		}
	}

	name := info.Name()

	// 检查名称匹配
	if len(rule.NameMatches) > 0 {
		matched := false
		for _, pattern := range rule.NameMatches {
			if s.matchPattern(name, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查文件扩展名
	if len(rule.FileExts) > 0 && !info.IsDir() {
		matched := false
		ext := strings.ToLower(filepath.Ext(name))
		for _, fileExt := range rule.FileExts {
			if strings.ToLower(fileExt) == ext {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查路径模式
	if len(rule.Patterns) > 0 {
		matched := false
		for _, pattern := range rule.Patterns {
			if s.matchPattern(path, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查基础路径
	if len(rule.BasePaths) > 0 {
		matched := false
		for _, basePath := range rule.BasePaths {
			expandedBase := fsutil.ExpandPath(basePath)
			if strings.HasPrefix(path, expandedBase) {
				matched = true
				break
			}
			// 支持通配符基础路径
			if s.matchPattern(path, basePath) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查最小大小
	if rule.MinSize > 0 && !info.IsDir() {
		if info.Size() < rule.MinSize {
			return false
		}
	}

	// 检查最大年龄
	if rule.MaxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -rule.MaxAge)
		if info.ModTime().After(cutoff) {
			return false
		}
	}

	return true
}

// matchPattern 匹配模式（支持通配符）
func (s *Scanner) matchPattern(str, pattern string) bool {
	// 处理路径展开
	str = fsutil.ExpandPath(str)
	pattern = fsutil.ExpandPath(pattern)

	// 简单的通配符匹配
	matched, err := filepath.Match(pattern, str)
	if err != nil {
		// 尝试作为简单字符串包含匹配
		return strings.Contains(str, pattern)
	}
	return matched
}
