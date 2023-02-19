package syscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/biz/cmdbiz"
)

// NewBatchRunCmd instance
func NewBatchRunCmd() *gcli.Command {
	var brOpts = struct {
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
			brOpts.BindCommonFlags(c)

			c.BoolOpt2(&brOpts.allSub, "all-subdir, all-sub", "run command on the each WORKDIR/subdir")
			c.VarOpt(&brOpts.exclude, "exclude", "e", "exclude some subdir on with --all-subdir")
			c.VarOpt(&brOpts.inDirs, "dirs", "", "run command on the each WORKDIR/dir, multi by comma")
			c.StrOpt2(&brOpts.cmdTpl, "cmd, c", "want execute command line, allow vars")
		},
		Func: func(c *gcli.Command, _ []string) error {

			return errors.New("TODO")
		},
	}
}
