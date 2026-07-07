package ghubcmd

import (
	"github.com/gookit/gcli/v3"
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

// NewGithubCmd commands
func NewGithubCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "github",
		Aliases: []string{"gh", "hub", "ghub"},
		Desc:    "useful tools for use github",
		Subs: []*gcli.Command{
			NewApiCmd(),
			NewDownloadAssetCmd(),
			gitcmd.NewPullRequestCmd(),
			gitcmd.NewBatchCmd(),
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

			c.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
				c.Infoln("[kite.GHub] Workdir:", c.WorkDir())
				return false
			})

			c.On(gcli.EvtCmdSubNotFound, gitcmd.RedirectToGitx)
		},
	}
}

func configProvider() *gitx.Config {
	return app.Ghub().Config
}

// NewDownloadAssetCmd instance
func NewDownloadAssetCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "down",
		Desc:    "checkout an new branch for development from `source` remote",
		Aliases: []string{"download"},
		Func: func(c *gcli.Command, args []string) error {

			return errorx.New("TODO")
		},
	}
}
