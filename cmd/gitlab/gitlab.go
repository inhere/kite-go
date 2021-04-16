package gitlab

import (
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/gituse"
)

var dryRun bool

var CmdForGitlab = &gcli.Command{
	Name: "gitlab",
	Aliases: []string{"gl", "gitl", "glab"},
	Desc: "useful tools for use gitlab",
	Subs: []*gcli.Command{
		UpdatePushCmd,
		UpdateNotPushCmd,
		gituse.OpenRemoteRepo,
	},
}

var upOpts = struct {
	notPush bool
}{}

var UpdatePushCmd = &gcli.Command{
	Name: "update-push",
	Aliases: []string{"up-push", "upp"},
	Desc: "Update codes from origin and main remote repositories, then push to remote",
	Config: func(c *gcli.Command) {
		UpdateNotPushCmd.Config(c)

		c.BoolVar(&upOpts.notPush, &gcli.FlagMeta{
			Name:  "not-push",
			Alias: "np",
			Desc:  "dont execute git push",
		})
	},
}

var UpdateNotPushCmd = &gcli.Command{
	Name: "update",
	Desc: "Update codes from origin and main remote repositories",
	Aliases: []string{"up"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&dryRun, "dry-run", "", false, "run workflow, but dont real execute command")

	},
}
