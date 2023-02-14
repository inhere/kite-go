package ghubcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/pkg/gitx"
)

// GithubCmd commands
var GithubCmd = &gcli.Command{
	Name:    "github",
	Aliases: []string{"gh", "hub", "ghub"},
	Desc:    "useful tools for use github",
	Subs: []*gcli.Command{
		gitcmd.NewAddCommitCmd(configProvider),
		gitcmd.NewAddCommitPush(configProvider),
		gitcmd.NewUpdateCmd(configProvider),
		gitcmd.NewUpdatePushCmd(configProvider),
		gitx.NewOpenRemoteCmd(func() string {
			return gitx.GitHubURL // github.host_url
		}),
	},
	Config: func(c *gcli.Command) {
		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			c.Infoln("[GHab] Workdir:", c.WorkDir())
			return false
		})

		c.On(events.OnCmdSubNotFound, gitcmd.RedirectToGitx)
	},
}

func configProvider() *gitx.Config {
	return app.Ghub().Config
}
