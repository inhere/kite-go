package gitx

import (
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/dump"
)

var bpullOpts = struct {
	dirs gcli.String
}{}

var BatchPull = &gcli.Command{
	Name:    "bpull",
	Desc:    "batch pull multi git directory by `git pull`",
	Aliases: []string{"bp", "bpul", "batch-pull"},
	Config: func(c *gcli.Command) {
		c.
			AddArg("baseDir", "base directory for run batch pull, default is work dir").
			With(func(arg *gcli.Argument) {
				arg.Value = "./"
			})

		c.VarOpt(&bpullOpts.dirs, "dirs", "", "limit update the given dir names")
	},
	Func: func(c *gcli.Command, args []string) error {
		baseDir := c.Arg("baseDir").String()
		absDir, err := filepath.Abs(baseDir)
		if err != nil {
			return err
		}

		dirNames := bpullOpts.dirs.Split(",")
		// if len(dirNames) > 0 {
		// 	for _, name := range dirNames {
		// 		path := filepath.Join(absDir, name)
		// 	}
		// 	return nil
		// }

		dump.P(dirNames, baseDir, absDir)

		return nil
	},
}

// UpdateCmd command
var UpdateCmd = &gcli.Command{
	Name:    "update",
	Desc:    "Update codes from git remote repositories",
	Aliases: []string{"up", "pul", "pull"},
	Config: func(c *gcli.Command) {
		bindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitw.NewWithArgs("pull", args...)
		pull.WithWorkDir(workdir)

		return pull.Run()
	},
}
