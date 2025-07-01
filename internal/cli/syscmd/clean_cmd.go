// clean_cmd provides command-line utilities for cleaning operations.
// This module includes various functions to assist with system commands.
package syscmd

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
)

type diskCleanOpts struct {
	ConfigFile  string   `flag:"desc=the config file path;name=config;shorts=c"`
	ScanDir     string   `flag:"desc=the start directory to scan;default=.;short=s"`
	Pattern     string   `flag:"desc=the pattern to match;short=p"`
	ExcludeDirs []string `flag:"desc=the exclude directories;name=exclude;short=E"`
	MaxDepth    int      `flag:"desc=the max depth to scan;default=-1;short=depth,D"`
	Concurrency int      `flag:"desc=the number of concurrent workers;default=3;short=C"`
	FileExts    []string `flag:"desc=the file extensions to match;short=e"`
	DryRun      bool     `flag:"desc=run in dry-run mode, do not delete files"`
}

// NewCleanCmd command
func NewCleanCmd() *gcli.Command {
	var cleanOpts = &diskCleanOpts{}

	return &gcli.Command{
		Name:    "clean",
		Desc:    "clean tmp or cache files or directories",
		Aliases: []string{"clear", "rm"},
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.FromStruct(cleanOpts))
		},
		Func: func(c *gcli.Command, _ []string) error {
			// do something
			return errorx.Raw("TODO")
		},
	}
}

// Config 配置结构体
type Config struct {
	ScanDir      string             `json:"scan_dir" yaml:"scan_dir"`
	Pattern      string             `json:"pattern" yaml:"pattern"`
	ExcludeDirs  []string           `json:"exclude_dirs" yaml:"exclude_dirs"`
	MaxDepth     int                `json:"max_depth" yaml:"max_depth"`
	Concurrency  int                `json:"concurrency" yaml:"concurrency"`
	FileExts     []string           `json:"file_exts" yaml:"file_exts"`
	DryRun       bool               `json:"dry_run" yaml:"dry_run"`
	RegexPattern *regexp.Regexp     `json:"-" yaml:"-"`
	OnMatch      func(string) error `json:"-" yaml:"-"`
}

type SysCleaner struct {
	cfg *Config
}

// 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		ScanDir:     ".",
		Pattern:     "*",
		ExcludeDirs: []string{},
		MaxDepth:    -1,
		Concurrency: 5,
		FileExts:    []string{},
		DryRun:      false,
	}

	// 支持 JSON/YAML 格式
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errorx.Wrap(err, "读取配置文件失败")
	}

	if strings.HasSuffix(configPath, ".json") {
		err = json.Unmarshal(data, config)
	} else if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
		err = yaml.Unmarshal(data, config)
	} else {
		return nil, errorx.Errorf("不支持的配置文件格式：%s", configPath)
	}

	// 验证配置
	if config.ScanDir == "" {
		return nil, errorx.New("起始目录不能为空")
	}
	if config.RegexPattern == nil && config.Pattern != "" {
		config.RegexPattern, err = regexp.Compile(config.Pattern)
		if err != nil {
			return nil, errorx.Wrap(err, "正则表达式模式编译失败")
		}
	}
	return config, nil
}

// 扩展通配符支持（兼容原有模式）
func matchWildcard(s, pattern string) bool {
	pattern = strings.ReplaceAll(pattern, `\*`, `.*`)
	pattern = strings.ReplaceAll(pattern, `\?`, `.`)
	return regexp.MustCompile(`^` + pattern + `$`).MatchString(s)
}
