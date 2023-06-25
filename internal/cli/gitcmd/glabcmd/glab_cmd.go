package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli/gitcmd"
	"github.com/inhere/kite-go/pkg/gitx"
)

var glOpts = struct {
	gitcmd.AutoChDir
}{}

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
		gitcmd.NewCheckoutCmd(),
		gitcmd.NewCloneCmd(configProvider),
		gitcmd.NewUpdateCmd(),
		gitcmd.NewUpdatePushCmd(),
		gitcmd.NewAddCommitPush(),
		gitcmd.NewAddCommitCmd(),
		gitcmd.NewOpenRemoteCmd(configProvider),
	},
	Config: func(c *gcli.Command) {
		glOpts.BindChdirFlags(c)

		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			wd := c.WorkDir()
			c.Infoln("[kite.GLAB] Workdir:", wd)
			return false
		})

		c.On(events.OnCmdSubNotFound, gitcmd.RedirectToGitx)
	},
}

func configProvider() *gitx.Config {
	return app.Glab().Config
}
