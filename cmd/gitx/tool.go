package gitx

import (
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

var (
	bpullOpts = struct {
		limit gcli.String
	}{}

	BatchPull = &gcli.Command{
		Name:    "bpull",
		Desc:    "batch pull multi git directory by `git pull`",
		Aliases: []string{"bp", "bpul", "batch-pull"},
		Config: func(c *gcli.Command) {
			c.AddArg("baseDir", "base directory for run batch pull, default is work dir").With(func(arg *gcli.Argument) {
				arg.Value = "./"
			})

			c.VarOpt(&bpullOpts.limit, "limit", "", "limit update the given dir names")
		},
		Func: func(c *gcli.Command, args []string) error {
			baseDir := c.Arg("baseDir").String()
			absDir, err := filepath.Abs(baseDir)
			if err != nil {
				return err
			}

			dump.P(bpullOpts.limit.Split(","), baseDir, absDir)

			return nil
		},
	}
)
