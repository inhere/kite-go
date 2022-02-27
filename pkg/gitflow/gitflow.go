package gitflow

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitwrap"
	"github.com/inherelab/kite/pkg/gituse"
)

var upOpts = struct {
	notPush bool
}{}

// UpdatePushCmd command
var UpdatePushCmd = &gcli.Command{
	Name:    "update-push",
	Desc:    "Update from origin and main remote, then push to origin remote",
	Aliases: []string{"up-push", "upp"},
	Config: func(c *gcli.Command) {
		gituse.BindCommonFlags(c)

		c.BoolVar(&upOpts.notPush, &gcli.FlagMeta{
			Name:  "not-push",
			Alias: "np",
			Desc:  "dont execute git push",
		})
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitwrap.Cmd("pull", args...)
		pull.WithWorkDir(gituse.Workdir)
		pull.OnBeforeExec(gitwrap.PrintCmdline)

		return pull.Run()
	},
}

// UpdateCmd command
var UpdateCmd = &gcli.Command{
	Name:    "update",
	Desc:    "Update from origin and main remote repositories",
	Aliases: []string{"up", "pul", "pull"},
	Config: func(c *gcli.Command) {
		gituse.BindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitwrap.Cmd("pull", args...)
		pull.WithWorkDir(gituse.Workdir)
		pull.OnBeforeExec(gitwrap.PrintCmdline)

		return pull.Run()
	},
}
