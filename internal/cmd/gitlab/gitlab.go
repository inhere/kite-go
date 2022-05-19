package gitlab

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/gitflow"
	"github.com/inherelab/kite/pkg/gituse"
)

var (
	dryRun  bool
	workdir string
)

func bindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&dryRun, "dry-run", "", false, "run workflow, but dont real execute command")
	c.StrOpt(&workdir, "workdir", "w", "", "the command workdir path")
}

// CmdForGitlab gitlab commands
var CmdForGitlab = &gcli.Command{
	Name:    "gitlab",
	Aliases: []string{"gl", "gitl", "glab"},
	Desc:    "useful tools for use gitlab",
	Subs: []*gcli.Command{
		gitflow.UpdateCmd,
		gitflow.UpdatePushCmd,
		gituse.OpenRemoteRepo,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdRunBefore, func(data ...interface{}) (stop bool) {
			color.Info.Println("Current workdir:", c.WorkDir())
			return false
		})
	},
}
