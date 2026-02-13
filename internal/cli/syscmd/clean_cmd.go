package syscmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/service/sysservice"
	"github.com/inhere/kite-go/internal/service/sysservice/sysclean"
)

type diskCleanOpts struct {
	// 配置
	ConfigFile string `flag:"desc=YAML/JSON配置文件路径;name=config;shorts=c"`

	// 扫描
	ScanDirs    gflag.Strings `flag:"desc=扫描目录;name=scan;short=s"`
	ExcludeDirs gflag.Strings `flag:"desc=排除目录;name=exclude;short=E"`
	MaxDepth    int           `flag:"desc=最大扫描深度，-1 表示无限制;default=-1;short=D"`
	Concurrency int           `flag:"desc=并发扫描数;default=5;short=C"`

	// 规则
	RuleNames  gflag.Strings `flag:"desc=指定规则名称;name=rule;short=R"`
	Categories gflag.Strings `flag:"desc=按类别过滤"`
	FileExts   gflag.Strings `flag:"desc=按文件扩展名过滤;short=e"`
	Pattern    string        `flag:"desc=自定义匹配模式;short=p"`

	// 行为
	DryRun   bool `flag:"desc=预览模式，不实际删除;short=dry"`
	UseTrash bool `flag:"desc=移动到回收站而非直接删除;name=trash"`
	Yes      bool `flag:"desc=跳过确认;short=y"`
	Force    bool `flag:"desc=强制执行高风险操作;short=f"`

	// 缓存
	UseCache   bool `flag:"desc=使用缓存;default=true"`
	ClearCache bool `flag:"desc=清除缓存"`

	// 输出
	OutputFormat string `flag:"desc=输出格式(text/markdown/json);default=markdown;short=o"`
	OutputFile   string `flag:"desc=输出文件路径;short=O"`
	Verbose      bool   `flag:"desc=详细输出;short=v"`
	ListRules    bool   `flag:"desc=列出所有可用规则"`
}

// NewSysCleanCmd 创建系统清理命令
func NewSysCleanCmd() *gcli.Command {
	var opts = &diskCleanOpts{}

	return &gcli.Command{
		Name:    "clean",
		Desc: "清理系统临时文件、缓存和无用数据",
		Aliases: []string{"clear", "rm"},
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.FromStruct(opts))
		},
		Func: func(c *gcli.Command, _ []string) error {
			return runCleanCmd(opts, c)
		},
	}
}

// runCleanCmd 执行清理命令
func runCleanCmd(opts *diskCleanOpts, c *gcli.Command) error {
	ctx := context.Background()

	// 构建配置
	cfg := sysclean.DefaultCleanConfig()
	cfg.ScanDirs = opts.ScanDirs
	cfg.ExcludeDirs = opts.ExcludeDirs
	cfg.MaxDepth = opts.MaxDepth
	cfg.Concurrency = opts.Concurrency
	cfg.RuleNames = opts.RuleNames
	cfg.FileExts = opts.FileExts
	cfg.Pattern = opts.Pattern
	cfg.DryRun = opts.DryRun
	cfg.UseTrash = opts.UseTrash
	cfg.Yes = opts.Yes
	cfg.Force = opts.Force
	cfg.UseCache = opts.UseCache
	cfg.OutputFormat = opts.OutputFormat
	cfg.OutputFile = opts.OutputFile
	cfg.Verbose = opts.Verbose

	// 解析类别
	if len(opts.Categories) > 0 {
		cfg.Categories = make([]sysclean.RuleCategory, 0, len(opts.Categories))
		for _, cat := range opts.Categories {
			cfg.Categories = append(cfg.Categories, sysclean.RuleCategory(cat))
		}
	}

	// 创建服务
	svc := sysservice.NewSysCleanService(cfg)

	// 加载配置文件
	if opts.ConfigFile != "" {
		if err := svc.LoadConfig(opts.ConfigFile); err != nil {
			return errorx.Wrap(err, "加载配置文件失败")
		}
	}

	// 列出规则
	if opts.ListRules {
		c.Println(svc.ListRules())
		return nil
	}

	// 清除缓存
	if opts.ClearCache {
		if err := svc.ClearCache(); err != nil {
			return errorx.Wrap(err, "清除缓存失败")
		}
		c.Println("缓存已清除")
		return nil
	}

	// 显示缓存信息
	if cfg.UseCache && cfg.Verbose {
		cacheInfo := svc.GetCacheInfo()
		if cacheInfo.Exists && !cacheInfo.IsExpired {
			c.Printf("使用缓存结果（剩余有效时间: %v）\n", cacheInfo.Remaining.Round(time.Second))
		}
	}

	// 扫描或清理
	var report *sysclean.CleanReport
	var err error

	if cfg.DryRun {
		// 预览模式
		c.Println("=== 预览模式（不会删除任何文件）===\n")
		report, err = svc.Preview(ctx)
	} else {
		// 执行清理
		if !cfg.Yes {
			// 先预览
			c.Println("=== 扫描结果预览 ===\n")
			report, err = svc.Preview(ctx)
			if err != nil {
				return err
			}

			// 显示摘要
			printScanSummary(report, c)

			// 确认
			if !confirmAction("是否继续执行清理?", c) {
				c.Println("操作已取消")
				return nil
			}
			c.Println()

			// 重新扫描并执行清理
			svc.ClearCache() // 清除缓存以重新扫描
			report, err = svc.Clean(ctx)
		} else {
			report, err = svc.Clean(ctx)
		}
	}

	if err != nil {
		return errorx.Wrap(err, "执行清理失败")
	}

	// 输出报告
	outputReport(svc, report, opts, c)

	return nil
}

// printScanSummary 打印扫描摘要
func printScanSummary(report *sysclean.CleanReport, c *gcli.Command) {
	if report.ScanStats == nil {
		return
	}

	c.Printf("扫描统计:\n")
	c.Printf("  - 扫描耗时: %v\n", report.ScanStats.Duration.Round(time.Millisecond))
	c.Printf("  - 匹配目标: %d 个\n", report.ScanStats.TotalTargets)
	c.Printf("  - 文件数: %d\n", report.ScanStats.TotalFiles)
	c.Printf("  - 目录数: %d\n", report.ScanStats.TotalDirs)
	c.Printf("  - 总大小: %s\n", formatSize(report.ScanStats.TotalSize))
	c.Println()

	// 显示高风险警告
	if len(report.Warnings) > 0 {
		c.Println("警告:")
		for _, warning := range report.Warnings {
			c.Printf("  ⚠️  %s\n", warning)
		}
		c.Println()
	}
}

// confirmAction 确认操作
func confirmAction(message string, c *gcli.Command) bool {
	c.Printf("%s [y/N]: ", message)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

// outputReport 输出报告
func outputReport(svc sysservice.SysCleanService, report *sysclean.CleanReport, opts *diskCleanOpts, c *gcli.Command) {
	var output string

	switch strings.ToLower(opts.OutputFormat) {
	case "text", "txt":
		output = svc.ExportMarkdown(report) // 暂时用 markdown
	case "json":
		output = "{\n  \"message\": \"JSON 格式暂未实现\"\n}"
	default:
		output = svc.ExportMarkdown(report)
	}

	// 输出到文件
	if opts.OutputFile != "" {
		if err := svc.SaveReport(report, opts.OutputFile); err != nil {
			c.Errorf("保存报告失败: %v\n", err)
		} else {
			c.Printf("报告已保存到: %s\n", opts.OutputFile)
		}
	}

	// 输出到控制台
	c.Println(output)
}

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
