package gitcmd

import (
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/inhere/kite/internal/biz/cmdbiz"
)

// BatchCmd command
var BatchCmd = &gcli.Command{
	Name:    "batch",
	Aliases: []string{"bat"},
	Desc:    "provide some useful dev tools commands",
	Subs: []*gcli.Command{
		BatchRunCmd,
		BatchStatusCmd,
		BatchPullCmd,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errorx.New("TODO")
	// },
}

var BatchStatusCmd = &gcli.Command{
	Name:    "status",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"st"},
}

var btrOpts = struct {
	cmdbiz.CommonOpts
	pDir string
}{}
var BatchRunCmd = &gcli.Command{
	Name:    "run",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"exec"},
	Config: func(c *gcli.Command) {
		btrOpts.BindCommonFlags(c)

		c.AddArg("dirs", "run command in the given dirs, if empty, run on all subdir")
	},
	Func: func(c *gcli.Command, args []string) error {

		return nil
	},
}

var btpOpts = struct {
	cmdbiz.CommonOpts
	dirs gcli.String
}{}

var BatchPullCmd = &gcli.Command{
	Name:    "pull",
	Desc:    "batch pull multi git directory by `git pull`",
	Aliases: []string{"pul", "pl"},
	Config: func(c *gcli.Command) {
		c.
			AddArg("baseDir", "base directory for run batch pull, default is work dir").
			WithValue("./")

		c.VarOpt(&btpOpts.dirs, "dirs", "", "limit update the given dir names")
	},
	Func: func(c *gcli.Command, args []string) error {
		baseDir := c.Arg("baseDir").String()
		absDir, err := filepath.Abs(baseDir)
		if err != nil {
			return err
		}

		dirNames := btpOpts.dirs.Split(",")
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