package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/inhere/kite/pkg/gitx"
)

var upOpts = struct {
	gitx.CommonOpts
}{}

// UpdateCmd command
var UpdateCmd = &gcli.Command{
	Name:    "pull",
	Desc:    "Update codes from git remote repositories",
	Aliases: []string{"pul", "pl"},
	Config: func(c *gcli.Command) {
		upOpts.BindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitw.NewWithArgs("pull", args...)
		pull.WithWorkDir(workdir)

		return pull.Run()
	},
}
