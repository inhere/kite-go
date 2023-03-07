package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
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
		ResolveConflictCmd,
		MergeRequestCmd,
		gitcmd.BatchCmd,
		gitcmd.NewInitFlowCmd(),
		gitcmd.NewBranchCmd(),
		gitcmd.NewCloneCmd(configProvider),
		gitcmd.NewUpdateCmd(configProvider),
		gitcmd.NewUpdatePushCmd(configProvider),
		gitcmd.NewAddCommitPush(configProvider),
		gitcmd.NewAddCommitCmd(configProvider),
		gitcmd.NewOpenRemoteCmd(configProvider),
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
