package fscmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
)

type dirTreeOpt struct {
	depth   uint
	exclude string
	include string
}

// NewDirTreeCmd 创建一个目录树命令
func NewDirTreeCmd() *gcli.Command {
	var dtOpt = dirTreeOpt{}

	return &gcli.Command{
		Name: "tree",
		Desc: "show directory tree",
		Config: func(cmd *gcli.Command) {
			cmd.UintOpt(&dtOpt.depth, "depth", "d", 10, "custom the tree max depth")
			cmd.StrOpt2(&dtOpt.exclude, "exclude,e", "custom the exclude name pattern")
			cmd.StrOpt2(&dtOpt.include, "include,i", "custom the include name pattern")

			cmd.AddArg("dir", "directory path").WithDefault(".")
		},
		Func: func(c *gcli.Command, args []string) error {
			dirPath := c.Arg("dir").String()
			return printDirTree(dirPath, dtOpt)
		},
	}
}

func printDirTree(root string, dtOpt dirTreeOpt) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		depth := uint(strings.Count(relPath, string(filepath.Separator)))
		if depth > dtOpt.depth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// is excluded path
		if fsutil.PathMatch(dtOpt.exclude, relPath) || !fsutil.PathMatch(dtOpt.include, relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		fmt.Printf("%s%s\n", strings.Repeat("  ", int(depth)), relPath)
		return nil
	})
}
