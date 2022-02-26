package gitlab

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/gituse"
)

var dryRun bool

// CmdForGitlab gitlab commands
var CmdForGitlab = &gcli.Command{
	Name:    "gitlab",
	Aliases: []string{"gl", "gitl", "glab"},
	Desc:    "useful tools for use gitlab",
	Subs: []*gcli.Command{
		UpdatePushCmd,
		UpdateNotPushCmd,
		gituse.OpenRemoteRepo,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdRunBefore, func(data ...interface{}) (stop bool) {
			color.Info.Println("Current workdir:", c.WorkDir())
			return false
		})
	},
}

var upOpts = struct {
	notPush bool
}{}

// UpdatePushCmd command
var UpdatePushCmd = &gcli.Command{
	Name:    "update-push",
	Desc:    "Update codes from origin and main remote repositories, then push to remote",
	Aliases: []string{"up-push", "upp"},
	Config: func(c *gcli.Command) {
		UpdateNotPushCmd.Config(c)

		c.BoolVar(&upOpts.notPush, &gcli.FlagMeta{
			Name:  "not-push",
			Alias: "np",
			Desc:  "dont execute git push",
		})
	},
	Func: func(c *gcli.Command, args []string) error {

		return nil
	},
}

// UpdateNotPushCmd command
var UpdateNotPushCmd = &gcli.Command{
	Name:    "update",
	Desc:    "Update codes from origin and main remote repositories",
	Aliases: []string{"up"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&dryRun, "dry-run", "", false, "run workflow, but dont real execute command")
	},
}
