package fscmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
)

// RenameOptions for rename command
type RenameOptions struct {
	Pattern     string `flag:"desc=regex pattern for matching files;shorts=p;required=true"`
	Replacement string `flag:"desc=replacement pattern;shorts=r;required=true"`
	DryRun      bool   `flag:"desc=dry run without actually renaming;shorts=d"`
	Verbose     bool   `flag:"desc=show verbose output;shorts=v"`
	Recursive   bool   `flag:"desc=recursive rename in sub-directories;shorts=R"`
	// Named       bool   `flag:"desc=use named var in replacement pattern"` // TODO

	// internal use
	dirs  []string
	paths []string
}

// NewRenameCmd instance
//
//	regex: from (\w+)_(\w+) to $1_new_$2
//	内置常见的匹配:
//	 - {word} 匹配单词 \w+
//	 - {num} 匹配数字 \d+
//	 - {name} 匹配名称 [\w-]+
//	 - {any} 匹配任意字符 .+
func NewRenameCmd() *gcli.Command {
	var renameOpts = &RenameOptions{}

	return &gcli.Command{
		Name: "rename",
		Desc: "batch rename files by regexp pattern",
		Help: `
<green>Built In Patterns</>:
> NOTE: N is index number, 1...N, on use multi times
 - {word[N]} matching word: \w+
 - {name[N]} matching name: [\w-]+
 - {numN} matching number: \d+
 - {alphaN} matching alpha: [a-zA-Z]+
 - {anyN} matching any: .+

<green>Built In Replacement Vars</>:
 - {ymd}  		  get current datetime: yyyy-mm-dd
 - {ymd_hms}  	  get current datetime: yyyy-mm-dd hh:mm:ss
 - {dirname}      get current dirname in filepath.
`,
		Examples: `
{$fullCmd} -p "(\w+)_(\w+)" -r "$1_new_$2" /path/to/directory
# or
{$fullCmd} -p "{word}_{word2}" -r "$1_new_$2" /path/to/directory

{$fullCmd} -p "(\d{4})-(\d{2})-(\d{2})" -r "$3-$2-$1" /path/to/files/*
# or
{$fullCmd} -p "{num}-{num2}-{num3}" -r "$3-$2-$1" /path/to/files/*
`,
		Config: func(c *gcli.Command) {
			c.MustFromStruct(renameOpts, gflag.TagRuleNamed)
			c.AddArg("files", "dir-path or files to rename, support glob pattern", true, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			if renameOpts.Pattern == "" || renameOpts.Replacement == "" {
				colorp.Errorln("Both pattern and replacement are required")
				return nil
			}

			// Get files to rename
			files := c.Arg("files").Strings()
			if len(files) == 0 {
				return fmt.Errorf("no file paths specified")
			}

			renameOpts.paths = files
			return handleRename(renameOpts)
		},
	}
}

var (
	regexPattern = regexp.MustCompile(`\{(\w+)\}`)
	regEndNumber = regexp.MustCompile(`\d+$`)

	builtInPatternVars = map[string]string{
		"word":  `(\w+)`,
		"num":   `(\d+)`,
		"name":  `([\w-]+)`,
		"any":   `(.+)`,
		"alpha": `([a-zA-Z]+)`,
	}
)

// 处理内置匹配,一个匹配多次使用需要加N 如 {word} {word1}
func formatRenamePattern(pattern string) string {
	return regexPattern.ReplaceAllStringFunc(pattern, func(match string) string {
		name := strings.Trim(match, "{}")
		ln := len(name)
		// 去除最后的N - 数字 1...N
		if name[ln-1] > '0' || name[ln-1] < '9' {
			if num := regEndNumber.FindString(name); num != "" {
				name = name[:ln-len(num)]
			}
		}

		if val, ok := builtInPatternVars[name]; ok {
			return val
		}
		return match
	})
}

func handleRename(opts *RenameOptions) error {
	// 处理内置匹配
	if strings.IndexByte(opts.Pattern, '{') >= 0 {
		opts.Pattern = formatRenamePattern(opts.Pattern)
	}

	// Compile regex
	colorp.Infof("Starting rename: %s -> %s\n", opts.Pattern, opts.Replacement)
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %v", err)
	}

	count := 0
	now := time.Now()
	var matches []string

	replStr := opts.Replacement
	replVars := map[string]string{
		"ymd":     now.Format("2006-01-02"),
		"ymd_hms": now.Format("2006-01-02 15:04:05"),
	}

	for _, globPath := range opts.paths {
		// 是一个目录路径,没有 glob 字符
		if !strutil.ContainsByteOne(globPath, []byte("*{")) && fsutil.IsDir(globPath) {
			globPath = globPath + "/*"
		}
		matches, err = filepath.Glob(globPath)
		if err != nil {
			colorp.Warnf("Error globbing %s: %v\n", globPath, err)
			continue
		}

		colorp.Infof("Processing %s, founds: %d\n", globPath, len(matches))
		for _, filePath := range matches {
			// Skip directories if not in recursive mode
			if !opts.Recursive && fsutil.IsDir(filePath) {
				continue
			}

			if opts.Verbose {
				fmt.Printf("[V] Found file: %s\n", filePath)
			}

			oldName := fsutil.Name(filePath)
			if strings.IndexByte(replStr, '{') >= 0 {
				replVars["dirname"] = fsutil.Name(fsutil.Dir(filePath))
				replStr = strutil.Replaces(replStr, replVars)
			}
			newName := re.ReplaceAllString(oldName, replStr)

			if newName != oldName {
				newPath := fsutil.JoinPaths(fsutil.Dir(filePath), newName)
				if opts.DryRun {
					count++
					colorp.Cyanf("- Dry-Run: would rename %s -> %s\n", oldName, newName)
					continue
				}

				colorp.Cyanf("Rename: %s -> %s\n", oldName, newName)
				// 检查 newPath 是否已经存在
				if fsutil.IsFile(newPath) {
					colorp.Warnf("New file already exists: %s\n", newPath)
					continue
				}

				if err = os.Rename(filePath, newPath); err != nil {
					colorp.Warnf("Failed to rename %s: %v\n", filePath, err)
				} else {
					count++
				}
			} else {
				colorp.Grayf("- Skipping %s: not-match or no-change\n", oldName)
			}
		}
	}

	if opts.DryRun {
		colorp.Infof("Dry run: would rename %d files\n", count)
	} else {
		colorp.Infof("Successfully renamed %d files\n", count)
	}

	return nil
}
