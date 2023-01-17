package gitlab

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/gitx"
	"github.com/inhere/kite/pkg/gitflow"
	"github.com/inhere/kite/pkg/gituse"
)

var (
	dryRun  bool
	workdir string
)

func bindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&dryRun, "dry-run", "dr", false, "run workflow, but dont real execute")
	c.StrOpt(&workdir, "workdir", "w", "", "custom the command workdir path")
}

// GitLab commands
var GitLab = &gcli.Command{
	Name:    "gitlab",
	Desc:    "useful tool commands for use gitlab",
	Aliases: []string{"gl", "glab"},
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
