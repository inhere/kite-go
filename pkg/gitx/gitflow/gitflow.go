package gitflow

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/inhere/kite/pkg/gitx"
)

var upOpts = struct {
	gitx.CommonOpts
	notPush bool
}{}

// UpdatePushCmd command
var UpdatePushCmd = &gcli.Command{
	Name:    "update-push",
	Desc:    "Update from origin and main remote, then push to origin remote",
	Aliases: []string{"up-push", "upp"},
	Config: func(c *gcli.Command) {
		upOpts.BindCommonFlags(c)

		c.BoolVar(&upOpts.notPush, &gcli.FlagMeta{
			Name:   "not-push",
			Desc:   "dont execute git push",
			Shorts: []string{"np"},
		})
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitw.Cmd("pull", args...)
		pull.WithWorkDir(upOpts.Workdir)
		pull.OnBeforeExec(gitw.PrintCmdline)

		return pull.Run()
	},
}

// UpdateCmd command
var UpdateCmd = &gcli.Command{
	Name:    "update",
	Desc:    "Update from origin and main remote repositories",
	Aliases: []string{"up", "pul", "pull"},
	Config: func(c *gcli.Command) {
		upOpts.BindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitw.Cmd("pull", args...)
		pull.WithWorkDir(upOpts.Workdir)
		pull.OnBeforeExec(gitw.PrintCmdline)

		return pull.Run()
	},
}
