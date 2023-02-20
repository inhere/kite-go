package syscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/biz/cmdbiz"
)

// NewBatchRunCmd instance
func NewBatchRunCmd() *gcli.Command {
	var btrOpts = struct {
		cmdbiz.CommonOpts
		cmdTpl  string
		inDirs  gcli.String
		allSub  bool
		exclude gcli.Strings
	}{}

	return &gcli.Command{
		Name:    "brun",
		Aliases: []string{"batch-run"},
		Desc:    "batch run more commands at once",
		Config: func(c *gcli.Command) {
			btrOpts.BindCommonFlags(c)

			c.BoolOpt2(&btrOpts.allSub, "all-subdir, all-sub", "run command on the each WORKDIR/subdir")
			c.VarOpt(&btrOpts.exclude, "exclude", "e", "exclude some subdir on with --all-subdir")
			c.VarOpt(&btrOpts.inDirs, "dirs", "", "run command on the each WORKDIR/dir, multi by comma")
			c.StrOpt2(&btrOpts.cmdTpl, "cmd, c", "want execute command line, allow vars")
		},
		Func: func(c *gcli.Command, _ []string) error {

			return errors.New("TODO")
		},
	}
}
