package sysclean

// GetPresetRules 获取预设规则
func GetPresetRules() []*CleanRule {
	rules := make([]*CleanRule, 0)

	// 通用规则
	rules = append(rules, getCommonRules()...)

	// 平台特定规则
	switch CurrentPlatform() {
	case PlatformWindows:
		rules = append(rules, getWindowsRules()...)
	case PlatformDarwin:
		rules = append(rules, getDarwinRules()...)
	case PlatformLinux:
		rules = append(rules, getLinuxRules()...)
	}

	return rules
}

// getCommonRules 获取通用规则（适用于所有平台）
func getCommonRules() []*CleanRule {
	return []*CleanRule{
		// Node.js 相关
		{
			Name:        "node_modules",
			Description: "Node.js 依赖目录",
			Category:    CategoryDependency,
			TargetType:  TargetTypeDir,
			NameMatches: []string{"node_modules"},
			RiskLevel:   1,
			Recursive:   false,
			Enabled:     true,
			ConfirmMsg:  "将删除 node_modules 目录（可通过 npm install 重新安装）",
		},
		{
			Name:        "npm_cache",
			Description: "npm 缓存目录",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/.npm"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "yarn_cache",
			Description: "Yarn 缓存目录",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/.yarn/cache", "~/Library/Caches/Yarn"},
			RiskLevel:   1,
			Enabled:     true,
		},

		// 日志文件
		{
			Name:        "log_files",
			Description: "日志文件（.log）",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			FileExts:    []string{".log", ".log.*"},
			RiskLevel:   1,
			Enabled:     true,
			ConfirmMsg:  "将删除 .log 日志文件",
		},
		{
			Name:        "npm_debug_log",
			Description: "npm 调试日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			NameMatches: []string{"npm-debug.log*", "yarn-debug.log*", "yarn-error.log*"},
			RiskLevel:   1,
			Enabled:     true,
		},

		// 临时文件
		{
			Name:        "temp_files",
			Description: "临时文件",
			Category:    CategoryTemp,
			TargetType:  TargetTypeFile,
			FileExts:    []string{".tmp", ".temp", ".bak", ".swp", ".swo"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "editor_backup",
			Description: "编辑器备份文件",
			Category:    CategoryTemp,
			TargetType:  TargetTypeFile,
			Patterns:    []string{"*~", "*.bak", "#*#"},
			RiskLevel:   1,
			Enabled:     true,
		},

		// Go 语言相关
		{
			Name:        "go_build_cache",
			Description: "Go 构建缓存",
			Category:    CategoryBuild,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/go/build", "~/Library/Caches/go-build"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "go_module_cache",
			Description: "Go 模块缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/go/pkg/mod"},
			RiskLevel:   2,
			Enabled:     true,
			ConfirmMsg:  "将删除 Go 模块缓存（重新构建时会自动下载）",
		},

		// Python 相关
		{
			Name:        "python_cache",
			Description: "Python __pycache__ 目录",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			NameMatches: []string{"__pycache__"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "python_pyc",
			Description: "Python 编译缓存文件",
			Category:    CategoryCache,
			TargetType:  TargetTypeFile,
			FileExts:    []string{".pyc", ".pyo", ".pyd"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "pip_cache",
			Description: "pip 缓存目录",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/.cache/pip", "~/Library/Caches/pip"},
			RiskLevel:   1,
			Enabled:     true,
		},

		// IDE 相关
		{
			Name:        "vscode_cache",
			Description: "VS Code 缓存目录",
			Category:    CategoryIDE,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{
				"~/.config/Code/Cache",
				"~/.config/Code/CachedData",
				"~/Library/Application Support/Code/Cache",
				"~/Library/Application Support/Code/CachedData",
				"~/AppData/Roaming/Code/Cache",
				"~/AppData/Roaming/Code/CachedData",
			},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "jetbrains_cache",
			Description: "JetBrains IDE 缓存",
			Category:    CategoryIDE,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{
				"~/.cache/JetBrains",
				"~/Library/Caches/JetBrains",
				"~/AppData/Local/JetBrains",
			},
			RiskLevel:   2,
			Enabled:     true,
			ConfirmMsg:  "将删除 JetBrains IDE 缓存（首次启动可能需要重建索引）",
		},

		// 版本控制
		{
			Name:        "git_orphan_objects",
			Description: "Git 孤立对象",
			Category:    CategorySystem,
			TargetType:  TargetTypeBoth,
			Patterns:    []string{".git/objects/*/tmp_*"},
			RiskLevel:   2,
			Enabled:     true,
		},

		// 构建产物
		{
			Name:        "dist_build",
			Description: "构建输出目录",
			Category:    CategoryBuild,
			TargetType:  TargetTypeDir,
			NameMatches: []string{"dist", "build", "out", "target"},
			RiskLevel:   2,
			Enabled:     true,
			ConfirmMsg:  "将删除 dist/build/out/target 目录（可通过重新构建生成）",
		},

		// 测试覆盖率
		{
			Name:        "test_coverage",
			Description: "测试覆盖率文件",
			Category:    CategoryTemp,
			TargetType:  TargetTypeFile,
			Patterns:    []string{"coverage.out", "coverage.html", "*.coverprofile"},
			RiskLevel:   1,
			Enabled:     true,
		},

		// Docker 相关
		{
			Name:        "docker_log",
			Description: "Docker 容器日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{
				"/var/lib/docker/containers/*/*-json.log",
			},
			RiskLevel:   2,
			MaxAge:      7,
			Enabled:     true,
		},
	}
}

// getWindowsRules 获取 Windows 平台特定规则
func getWindowsRules() []*CleanRule {
	return []*CleanRule{
		{
			Name:        "windows_temp",
			Description: "Windows 临时文件夹",
			Category:    CategoryTemp,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{
				"~/AppData/Local/Temp",
				"/Windows/Temp",
			},
			RiskLevel:   1,
			MaxAge:      7,
			Enabled:     true,
		},
		{
			Name:        "windows_update_cache",
			Description: "Windows 更新缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{
				"/Windows/SoftwareDistribution/Download",
			},
			RiskLevel:   2,
			Enabled:     true,
			ConfirmMsg:  "将删除 Windows 更新缓存文件",
		},
		{
			Name:        "windows_thumbnail_cache",
			Description: "Windows 缩略图缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{"~/AppData/Local/Microsoft/Windows/Explorer"},
			Patterns:    []string{"thumbcache_*.db"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "windows_prefetch",
			Description: "Windows 预读取文件",
			Category:    CategorySystem,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{"/Windows/Prefetch"},
			FileExts:    []string{".pf"},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "windows_recycle_bin",
			Description: "回收站（已删除但未清空）",
			Category:    CategorySystem,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"$RECYCLE.BIN"},
			RiskLevel:   3,
			Enabled:     true,
			ConfirmMsg:  "将清空回收站",
		},
		{
			Name:        "browser_cache_windows",
			Description: "浏览器缓存（Chrome/Edge/Firefox）",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths: []string{
				"~/AppData/Local/Google/Chrome/User Data/Default/Cache",
				"~/AppData/Local/Microsoft/Edge/User Data/Default/Cache",
				"~/AppData/Local/Mozilla/Firefox/Profiles/*/cache2",
			},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "windows_error_reports",
			Description: "Windows 错误报告",
			Category:    CategoryLog,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/AppData/Local/Microsoft/Windows/WER"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "windows_logs",
			Description: "Windows 日志文件",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{"/Windows/Logs"},
			FileExts:    []string{".log", ".etl"},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "windows_old",
			Description: "Windows.old 文件夹（系统升级残留）",
			Category:    CategorySystem,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"/Windows.old"},
			RiskLevel:   3,
			Enabled:     true,
			ConfirmMsg:  "将删除 Windows.old 文件夹（无法回退到旧版本 Windows）",
		},
	}
}

// getDarwinRules 获取 macOS 平台特定规则
func getDarwinRules() []*CleanRule {
	return []*CleanRule{
		{
			Name:        "macos_cache",
			Description: "macOS 系统缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/Library/Caches"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "macos_logs",
			Description: "macOS 系统日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/Library/Logs", "/var/log"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "macos_user_logs",
			Description: "macOS 用户日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/Library/Logs"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "xcode_derived_data",
			Description: "Xcode DerivedData",
			Category:    CategoryBuild,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Developer/Xcode/DerivedData"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "xcode_archives",
			Description: "Xcode 归档文件",
			Category:    CategoryBuild,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Developer/Xcode/Archives"},
			RiskLevel:   2,
			Enabled:     true,
			ConfirmMsg:  "将删除 Xcode 归档文件（建议先备份）",
		},
		{
			Name:        "xcode_device_support",
			Description: "Xcode 设备支持文件",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Developer/Xcode/iOS DeviceSupport"},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "homebrew_cache",
			Description: "Homebrew 缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Caches/Homebrew"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "homebrew_logs",
			Description: "Homebrew 日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{"~/Library/Logs/Homebrew"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "dmg_residual",
			Description: "DMG 安装包残留",
			Category:    CategoryTemp,
			TargetType:  TargetTypeFile,
			FileExts:    []string{".dmg"},
			Patterns:    []string{"~/Downloads/*.dmg"},
			MaxAge:      30,
			RiskLevel:   1,
			Enabled:     true,
			ConfirmMsg:  "将删除下载目录中超过 30 天的 DMG 文件",
		},
		{
			Name:        "macos_trash",
			Description: "macOS 废纸篓",
			Category:    CategorySystem,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/.Trash"},
			RiskLevel:   3,
			Enabled:     true,
			ConfirmMsg:  "将清空废纸篓",
		},
		{
			Name:        "macos_application_support",
			Description: "已卸载应用的残留数据",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Application Support"},
			RiskLevel:   2,
			Enabled:     false, // 默认禁用，需要手动检查
		},
		{
			Name:        "cocoapods_cache",
			Description: "CocoaPods 缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Caches/CocoaPods"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "carthage_cache",
			Description: "Carthage 缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/Library/Caches/org.carthage.CarthageKit"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "swift_package_cache",
			Description: "Swift Package Manager 缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{
				"~/Library/Caches/org.swift.swiftpm",
				"~/Library/Developer/Xcode/DerivedData/*/SourcePackages",
			},
			RiskLevel:   1,
			Enabled:     true,
		},
	}
}

// getLinuxRules 获取 Linux 平台特定规则
func getLinuxRules() []*CleanRule {
	return []*CleanRule{
		{
			Name:        "linux_temp",
			Description: "Linux 临时文件",
			Category:    CategoryTemp,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/tmp", "/var/tmp"},
			RiskLevel:   1,
			MaxAge:      7,
			Enabled:     true,
		},
		{
			Name:        "apt_cache",
			Description: "APT 包管理器缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/var/cache/apt/archives"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "dnf_cache",
			Description: "DNF/YUM 包管理器缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/var/cache/dnf", "/var/cache/yum"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "user_cache_linux",
			Description: "用户缓存目录",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/.cache"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "system_logs_linux",
			Description: "系统日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeFile,
			BasePaths:   []string{"/var/log"},
			FileExts:    []string{".log", ".gz", ".old"},
			RiskLevel:   2,
			MaxAge:      30,
			Enabled:     true,
		},
		{
			Name:        "journal_logs",
			Description: "systemd journal 日志",
			Category:    CategoryLog,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/var/log/journal"},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "package_cache_linux",
			Description: "包管理器缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{
				"/var/cache/pacman/pkg",
				"/var/cache/xbps",
			},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "linux_thumbnail_cache",
			Description: "Linux 缩略图缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeDir,
			BasePaths:   []string{"~/.cache/thumbnails"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "linux_trash",
			Description: "Linux 回收站",
			Category:    CategorySystem,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/.local/share/Trash"},
			RiskLevel:   3,
			Enabled:     true,
			ConfirmMsg:  "将清空回收站",
		},
		{
			Name:        "snap_cache",
			Description: "Snap 包缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{
				"/var/lib/snapd/cache",
				"~/snap",
			},
			RiskLevel:   2,
			Enabled:     true,
		},
		{
			Name:        "flatpak_cache",
			Description: "Flatpak 缓存",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"~/.local/share/flatpak/repo/tmp"},
			RiskLevel:   1,
			Enabled:     true,
		},
		{
			Name:        "old_kernels",
			Description: "旧内核文件",
			Category:    CategorySystem,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/boot"},
			Patterns:    []string{"initrd.img-*", "vmlinuz-*", "System.map-*"},
			RiskLevel:   3,
			Enabled:     false, // 危险操作，默认禁用
			ConfirmMsg:  "将删除旧内核文件（请确保当前内核正常工作）",
		},
		{
			Name:        "docker_linux_cache",
			Description: "Docker 缓存和未使用镜像",
			Category:    CategoryCache,
			TargetType:  TargetTypeBoth,
			BasePaths:   []string{"/var/lib/docker"},
			RiskLevel:   2,
			Enabled:     false, // 需要特殊处理
		},
	}
}

// FilterRulesByNames 按名称过滤规则
func FilterRulesByNames(rules []*CleanRule, names []string) []*CleanRule {
	if len(names) == 0 {
		return rules
	}

	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	result := make([]*CleanRule, 0)
	for _, rule := range rules {
		if nameSet[rule.Name] && rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// FilterRulesByCategories 按类别过滤规则
func FilterRulesByCategories(rules []*CleanRule, categories []RuleCategory) []*CleanRule {
	if len(categories) == 0 {
		return rules
	}

	catSet := make(map[RuleCategory]bool)
	for _, cat := range categories {
		catSet[cat] = true
	}

	result := make([]*CleanRule, 0)
	for _, rule := range rules {
		if catSet[rule.Category] && rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// FilterRulesByRiskLevel 按风险等级过滤规则
func FilterRulesByRiskLevel(rules []*CleanRule, maxRisk int) []*CleanRule {
	if maxRisk <= 0 {
		maxRisk = 2 // 默认只允许低和中风险
	}

	result := make([]*CleanRule, 0)
	for _, rule := range rules {
		if rule.RiskLevel <= maxRisk && rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// FilterEnabledRules 过滤启用的规则
func FilterEnabledRules(rules []*CleanRule) []*CleanRule {
	result := make([]*CleanRule, 0)
	for _, rule := range rules {
		if rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// FilterRulesByPlatform 过滤适用于当前平台的规则
func FilterRulesByPlatform(rules []*CleanRule, platform Platform) []*CleanRule {
	result := make([]*CleanRule, 0)
	for _, rule := range rules {
		// 如果没有指定平台，则适用于所有平台
		if len(rule.Platforms) == 0 {
			result = append(result, rule)
			continue
		}
		// 检查是否适用于当前平台
		for _, p := range rule.Platforms {
			if p == platform || p == PlatformAll {
				result = append(result, rule)
				break
			}
		}
	}
	return result
}
