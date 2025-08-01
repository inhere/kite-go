package ghubcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/internal/cli/gitcmd"
	"github.com/inhere/kite-go/pkg/gitx"
)

// GhOpts object
var GhOpts = struct {
	gitcmd.AutoGitDir
	cmdbiz.CommonOpts
}{}

// GithubCmd commands
var GithubCmd = &gcli.Command{
	Name:    "github",
	Aliases: []string{"gh", "hub", "ghub"},
	Desc:    "useful tools for use github",
	Subs: []*gcli.Command{
		DownloadAssetCmd,
		gitcmd.NewPullRequestCmd(),
		gitcmd.BatchCmd,
		gitcmd.NewBranchCmd(),
		gitcmd.NewCloneCmd(configProvider),
		gitcmd.NewAddCommitCmd(),
		gitcmd.NewAddCommitPush(),
		gitcmd.NewUpdateCmd(),
		gitcmd.NewUpdatePushCmd(),
		gitcmd.NewOpenRemoteCmd(configProvider),
	},
	Config: func(c *gcli.Command) {
		GhOpts.BindCommonFlags(c)
		GhOpts.BindChdirFlags(c)

		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			c.Infoln("[kite.GHab] Workdir:", c.WorkDir())
			return false
		})

		c.On(events.OnCmdSubNotFound, gitcmd.RedirectToGitx)
	},
}

func configProvider() *gitx.Config {
	return app.Ghub().Config
}

// DownloadAssetCmd instance
var DownloadAssetCmd = &gcli.Command{
	Name:    "down",
	Desc:    "checkout an new branch for development from `source` remote",
	Aliases: []string{"download"},
	Func: func(c *gcli.Command, args []string) error {

		return errorx.New("TODO")
	},
}
