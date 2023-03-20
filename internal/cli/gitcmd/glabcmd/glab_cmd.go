package glabcmd

import (
	"os"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli/gitcmd"
	"github.com/inhere/kite-go/pkg/gitx"
)

var glOpts = struct {
	autoGit bool
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
		gitcmd.NewCloneCmd(configProvider),
		gitcmd.NewUpdateCmd(configProvider),
		gitcmd.NewUpdatePushCmd(configProvider),
		gitcmd.NewAddCommitPush(configProvider),
		gitcmd.NewAddCommitCmd(configProvider),
		gitcmd.NewOpenRemoteCmd(configProvider),
	},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&glOpts.autoGit, "auto-root, auto-git", "auto change workdir to git repo root dir")

		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			wd := c.WorkDir()
			c.Infoln("[GLab] Workdir:", wd)

			if glOpts.autoGit {
				// auto find .git dir in parent.
				repoDir, changed := fsutil.SearchNameUpx(wd, gitw.GitDir)
				if changed {
					goutil.MustOK(os.Chdir(repoDir))
					c.ChWorkDir(repoDir)
					cliutil.Yellowf("NOTICE: auto founded git root and will chdir to: %s\n", repoDir)
				}
			}

			return false
		})

		c.On(events.OnCmdSubNotFound, gitcmd.RedirectToGitx)
	},
}

func configProvider() *gitx.Config {
	return app.Glab().Config
}
