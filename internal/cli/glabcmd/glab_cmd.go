package glabcmd

import (
	"os"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/gitflow"
)

// GitLabCmd commands
var GitLabCmd = &gcli.Command{
	Name:    "gitlab",
	Desc:    "useful tool commands for use gitlab",
	Aliases: []string{"gl", "glab"},
	Subs: []*gcli.Command{
		gitflow.UpdateCmd,
		gitflow.UpdatePushCmd,
		// gitx.NewOpenRemoteCmd(app.Cfg().String("gitlab.host_url")),
		gitx.NewOpenRemoteCmd(os.Getenv("KITE_GLAB_HOST")),
		MergeRequestCmd,
		gitcmd.AddCommitPush,
		gitcmd.AddCommitNotPush,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
			color.Info.Println("Current workdir:", c.WorkDir())
			return false
		})
	},
}
