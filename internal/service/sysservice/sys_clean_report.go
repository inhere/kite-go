package sysservice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/jsonutil"
	"github.com/inhere/kite-go/internal/service/sysservice/sysclean"
)

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// GenerateReport 生成清理报告
func (s *sysCleanService) GenerateReport(result *sysclean.ScanResult, isExecuted bool) *sysclean.CleanReport {
	report := &sysclean.CleanReport{
		GeneratedAt: time.Now(),
		Platform:    result.Platform,
		Mode:        "dry-run",
		ByCategory:  make(map[sysclean.RuleCategory]*sysclean.CategoryStats),
		ByRule:      make(map[string]*sysclean.RuleStats),
		Warnings:    make([]string, 0),
		Errors:      make([]string, 0),
	}

	// 获取主机名和用户名
	if hostname, err := os.Hostname(); err == nil {
		report.Hostname = hostname
	}
	if currentUser, err := user.Current(); err == nil {
		report.Username = currentUser.Username
	}

	if isExecuted {
		report.Mode = "clean"
	}

	// 扫描统计
	report.ScanStats = &sysclean.ScanStats{
		TotalFiles:   result.TotalFiles,
		TotalDirs:    result.TotalDirs,
		TotalTargets: result.TotalTargets,
		TotalSize:    result.TotalSize,
		Duration:     result.ScanDuration,
	}

	// 按类别和规则统计
	for _, group := range result.Groups {
		// 按类别统计
		catStats, ok := report.ByCategory[group.Category]
		if !ok {
			catStats = &sysclean.CategoryStats{Category: group.Category}
			report.ByCategory[group.Category] = catStats
		}
		catStats.FileCount += group.FileCount
		catStats.DirCount += group.DirCount
		catStats.TotalSize += group.TotalSize

		// 按规则统计
		report.ByRule[group.RuleName] = &sysclean.RuleStats{
			RuleName:  group.RuleName,
			FileCount: group.FileCount,
			DirCount:  group.DirCount,
			TotalSize: group.TotalSize,
		}
	}

	// 添加错误到报告
	for _, err := range result.Errors {
		report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", err.Path, err.Error))
	}

	// 添加风险警告
	for _, group := range result.Groups {
		if group.RiskLevel >= 3 {
			report.Warnings = append(report.Warnings,
				fmt.Sprintf("规则 '%s' 为高风险操作，请谨慎处理", group.RuleName))
		}
	}

	return report
}

// ExportMarkdown 导出 Markdown 格式报告
func (s *sysCleanService) ExportMarkdown(report *sysclean.CleanReport) string {
	var buf bytes.Buffer

	// 标题
	buf.WriteString("# 系统清理报告\n\n")

	// 基本信息
	buf.WriteString(fmt.Sprintf("**生成时间**: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("**平台**: %s\n", report.Platform))
	if report.Hostname != "" {
		buf.WriteString(fmt.Sprintf("**主机名**: %s\n", report.Hostname))
	}
	if report.Username != "" {
		buf.WriteString(fmt.Sprintf("**用户**: %s\n", report.Username))
	}
	buf.WriteString(fmt.Sprintf("**模式**: %s\n\n", report.Mode))

	buf.WriteString("---\n\n")

	// 执行摘要
	buf.WriteString("## 执行摘要\n\n")
	buf.WriteString("| 指标 | 值 |\n")
	buf.WriteString("|------|----|\n")

	if report.ScanStats != nil {
		buf.WriteString(fmt.Sprintf("| 扫描耗时 | %v |\n", report.ScanStats.Duration.Round(time.Millisecond)))
		buf.WriteString(fmt.Sprintf("| 匹配目标数 | %d |\n", report.ScanStats.TotalTargets))
		buf.WriteString(fmt.Sprintf("| 扫描总大小 | %s |\n", formatSize(report.ScanStats.TotalSize)))
	}

	if report.CleanStats != nil {
		buf.WriteString(fmt.Sprintf("| 释放空间 | %s |\n", formatSize(report.CleanStats.FreedSpace)))
		buf.WriteString(fmt.Sprintf("| 删除文件数 | %d |\n", report.CleanStats.DeletedFiles))
		buf.WriteString(fmt.Sprintf("| 删除目录数 | %d |\n", report.CleanStats.DeletedDirs))
		if report.CleanStats.FailedCount > 0 {
			buf.WriteString(fmt.Sprintf("| 失败数 | %d |\n", report.CleanStats.FailedCount))
		}
	}
	buf.WriteString("\n")

	// 分类统计
	if len(report.ByCategory) > 0 {
		buf.WriteString("## 分类统计\n\n")
		buf.WriteString("| 类别 | 文件数 | 目录数 | 总大小 |\n")
		buf.WriteString("|------|--------|--------|--------|\n")

		// 排序类别
		categories := make([]sysclean.RuleCategory, 0, len(report.ByCategory))
		for cat := range report.ByCategory {
			categories = append(categories, cat)
		}
		sort.Slice(categories, func(i, j int) bool {
			return report.ByCategory[categories[i]].TotalSize > report.ByCategory[categories[j]].TotalSize
		})

		for _, cat := range categories {
			stats := report.ByCategory[cat]
			buf.WriteString(fmt.Sprintf("| %s | %d | %d | %s |\n",
				cat.DisplayName(),
				stats.FileCount,
				stats.DirCount,
				formatSize(stats.TotalSize)))
		}
		buf.WriteString("\n")
	}

	// 规则详情（按风险等级分组）
	if len(report.ByRule) > 0 {
		buf.WriteString("## 规则详情\n\n")

		// 按风险等级分组
		riskGroups := make(map[int][]string)
		for ruleName := range report.ByRule {
			rule, found := s.GetRule(ruleName)
			riskLevel := 1
			if found {
				riskLevel = rule.RiskLevel
			}
			riskGroups[riskLevel] = append(riskGroups[riskLevel], ruleName)
		}

		// 按风险等级从低到高输出
		for risk := 1; risk <= 3; risk++ {
			ruleNames, ok := riskGroups[risk]
			if !ok || len(ruleNames) == 0 {
				continue
			}

			buf.WriteString(fmt.Sprintf("### 风险等级 %d\n\n", risk))
			buf.WriteString("| 规则名称 | 文件数 | 目录数 | 大小 |\n")
			buf.WriteString("|----------|--------|--------|------|\n")

			// 按大小排序
			sort.Slice(ruleNames, func(i, j int) bool {
				return report.ByRule[ruleNames[i]].TotalSize > report.ByRule[ruleNames[j]].TotalSize
			})

			for _, ruleName := range ruleNames {
				stats := report.ByRule[ruleName]
				buf.WriteString(fmt.Sprintf("| %s | %d | %d | %s |\n",
					ruleName,
					stats.FileCount,
					stats.DirCount,
					formatSize(stats.TotalSize)))
			}
			buf.WriteString("\n")
		}
	}

	// 警告
	if len(report.Warnings) > 0 {
		buf.WriteString("## 警告\n\n")
		for _, warning := range report.Warnings {
			buf.WriteString(fmt.Sprintf("- ⚠️ %s\n", warning))
		}
		buf.WriteString("\n")
	}

	// 错误
	if len(report.Errors) > 0 {
		buf.WriteString("## 错误\n\n")
		for _, err := range report.Errors {
			buf.WriteString(fmt.Sprintf("- ❌ %s\n", err))
		}
		buf.WriteString("\n")
	}

	// 页脚
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("\n*由 kite-go sys clean 生成于 %s*\n", time.Now().Format("2006-01-02 15:04:05")))

	return buf.String()
}

// ExportJSON 导出 JSON 格式报告
func (s *sysCleanService) ExportJSON(report *sysclean.CleanReport) (string, error) {
	return jsonutil.EncodeString( report)
}

// ExportText 导出文本格式报告
func (s *sysCleanService) ExportText(report *sysclean.CleanReport) string {
	var buf bytes.Buffer

	buf.WriteString("=== 系统清理报告 ===\n\n")

	// 基本信息
	buf.WriteString(fmt.Sprintf("生成时间: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("平台: %s\n", report.Platform))
	buf.WriteString(fmt.Sprintf("模式: %s\n\n", report.Mode))

	// 统计摘要
	if report.ScanStats != nil {
		buf.WriteString("--- 扫描统计 ---\n")
		buf.WriteString(fmt.Sprintf("扫描耗时: %v\n", report.ScanStats.Duration.Round(time.Millisecond)))
		buf.WriteString(fmt.Sprintf("匹配目标: %d\n", report.ScanStats.TotalTargets))
		buf.WriteString(fmt.Sprintf("总大小: %s\n\n", formatSize(report.ScanStats.TotalSize)))
	}

	if report.CleanStats != nil {
		buf.WriteString("--- 清理统计 ---\n")
		buf.WriteString(fmt.Sprintf("释放空间: %s\n", formatSize(report.CleanStats.FreedSpace)))
		buf.WriteString(fmt.Sprintf("删除文件: %d\n", report.CleanStats.DeletedFiles))
		buf.WriteString(fmt.Sprintf("删除目录: %d\n\n", report.CleanStats.DeletedDirs))
	}

	// 分类统计
	if len(report.ByCategory) > 0 {
		buf.WriteString("--- 分类统计 ---\n")
		for cat, stats := range report.ByCategory {
			buf.WriteString(fmt.Sprintf("%s: %d 文件, %d 目录, %s\n",
				cat.DisplayName(),
				stats.FileCount,
				stats.DirCount,
				formatSize(stats.TotalSize)))
		}
		buf.WriteString("\n")
	}

	// 规则详情
	if len(report.ByRule) > 0 {
		buf.WriteString("--- 规则详情 ---\n")
		for ruleName, stats := range report.ByRule {
			buf.WriteString(fmt.Sprintf("%s: %d 文件, %d 目录, %s\n",
				ruleName,
				stats.FileCount,
				stats.DirCount,
				formatSize(stats.TotalSize)))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// SaveReport 保存报告到文件
func (s *sysCleanService) SaveReport(report *sysclean.CleanReport, filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errorx.Wrap(err, "创建报告目录失败")
	}

	var content string
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".md", ".markdown":
		content = s.ExportMarkdown(report)
	case ".json":
		jsonContent, err := s.ExportJSON(report)
		if err != nil {
			return err
		}
		content = jsonContent
	default:
		content = s.ExportText(report)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return errorx.Wrap(err, "写入报告文件失败")
	}

	return nil
}

// Preview 预览模式（生成报告但不执行清理）
func (s *sysCleanService) Preview(ctx context.Context) (*sysclean.CleanReport, error) {
	result, err := s.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return s.GenerateReport(result, false), nil
}

// Clean 执行清理
func (s *sysCleanService) Clean(ctx context.Context) (*sysclean.CleanReport, error) {
	// 扫描
	result, err := s.Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 生成报告
	report := s.GenerateReport(result, true)

	// 如果是 dry-run 模式，只返回报告不执行清理
	if s.config.DryRun {
		report.Mode = "dry-run"
		return report, nil
	}

	// 执行清理
	cleanStats := &sysclean.CleanStats{
		UseTrash: s.config.UseTrash,
	}
	startTime := time.Now()

	for _, group := range result.Groups {
		// 检查风险等级
		if group.RiskLevel >= 3 && !s.config.Force {
			report.Warnings = append(report.Warnings,
				fmt.Sprintf("跳过高风险规则 '%s'（使用 --force 强制执行）", group.RuleName))
			continue
		}

		for _, target := range group.Targets {
			err := sysclean.DeleteTarget(target.Path, s.config.UseTrash, s.trashManager)
			if err != nil {
				cleanStats.FailedCount++
				cleanStats.FailedSpace += target.Size
				report.Errors = append(report.Errors,
					fmt.Sprintf("清理失败 %s: %s", target.Path, err.Error()))
			} else {
				if target.IsDir {
					cleanStats.DeletedDirs++
				} else {
					cleanStats.DeletedFiles++
				}
				cleanStats.FreedSpace += target.Size
			}
		}
	}

	cleanStats.Duration = time.Since(startTime)
	report.CleanStats = cleanStats

	// 清除缓存
	_ = s.ClearCache()

	return report, nil
}

// ListRules 列出所有可用规则
func (s *sysCleanService) ListRules() string {
	var buf bytes.Buffer

	buf.WriteString("=== 可用清理规则 ===\n\n")

	// 按类别分组
	categories := make(map[sysclean.RuleCategory][]*sysclean.CleanRule)
	for _, rule := range s.presetRules {
		categories[rule.Category] = append(categories[rule.Category], rule)
	}

	// 类别顺序
	catOrder := []sysclean.RuleCategory{
		sysclean.CategoryCache, sysclean.CategoryLog, sysclean.CategoryTemp,
		sysclean.CategoryBuild, sysclean.CategoryDependency, sysclean.CategoryIDE,
		sysclean.CategorySystem, sysclean.CategoryCustom,
	}

	for _, cat := range catOrder {
		rules := categories[cat]
		if len(rules) == 0 {
			continue
		}

		buf.WriteString(fmt.Sprintf("## %s\n", cat.DisplayName()))

		for _, rule := range rules {
			status := "✓"
			if !rule.Enabled {
				status = "○"
			}
			riskStars := strings.Repeat("★", rule.RiskLevel) + strings.Repeat("☆", 3-rule.RiskLevel)
			buf.WriteString(fmt.Sprintf("  %s %-20s [%s] %s\n",
				status, rule.Name, riskStars, rule.Description))
		}
		buf.WriteString("\n")
	}

	buf.WriteString("图例: ✓=启用 ○=禁用 ★=风险等级(1-3)\n")

	return buf.String()
}
