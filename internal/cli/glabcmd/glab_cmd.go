package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/envutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/pkg/gitx"
)

// GitLabCmd commands
var GitLabCmd = &gcli.Command{
	Name:    "gitlab",
	Desc:    "useful tool commands for use gitlab",
	Aliases: []string{"gl", "glab"},
	Subs: []*gcli.Command{
		MergeRequestCmd,
		gitcmd.NewUpdateCmd(configProvider),
		gitcmd.NewUpdatePushCmd(configProvider),
		gitcmd.NewAddCommitPush(configProvider),
		gitcmd.NewAddCommitCmd(configProvider),
		gitx.NewOpenRemoteCmd(func() string {
			return envutil.Getenv("KITE_GLAB_HOST", app.Cfg().String("gitlab.host_url"))
		}),
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			c.Infoln("[GLab] Workdir:", c.WorkDir())
			return false
		})

		c.On(events.OnCmdSubNotFound, gitcmd.RedirectToGitx)
	},
}

func configProvider() *gitx.Config {
	return app.Glab().Config
}
