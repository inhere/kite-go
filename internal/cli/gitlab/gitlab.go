package gitlab

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/internal/cli/gitx"
	"github.com/inherelab/kite/pkg/gitflow"
	"github.com/inherelab/kite/pkg/gituse"
)

var (
	dryRun  bool
	workdir string
)

func bindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&dryRun, "dry-run", "dr", false, "run workflow, but dont real execute command")
	c.StrOpt(&workdir, "workdir", "w", "", "the command workdir path")
}

// GitLab commands
var GitLab = &gcli.Command{
	Name:    "gitlab",
	Desc:    "useful tool commands for use gitlab",
	Aliases: []string{"gl", "gitl", "glab"},
	Subs: []*gcli.Command{
		gitflow.UpdateCmd,
		gitflow.UpdatePushCmd,
		gituse.OpenRemoteRepo,
		MergeRequest,
		gitx.AddCommitPush,
		gitx.AddCommitNotPush,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			color.Info.Println("Current workdir:", c.WorkDir())
			return false
		})
	},
}
