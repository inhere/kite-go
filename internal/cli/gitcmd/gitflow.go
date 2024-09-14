package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

// NewInitFlowCmd instance
func NewInitFlowCmd() *gcli.Command {
	var ifOpts = struct {
		cmdbiz.CommonOpts
	}{}

	return &gcli.Command{
		Name:    "init",
		Desc:    "init repo remote and other info for current project",
		Aliases: []string{"setup"},
		Config: func(c *gcli.Command) {
			ifOpts.BindWorkdirDryRun(c)
		},
		Func: func(c *gcli.Command, args []string) (err error) {
			cfg := apputil.GitCfgByCmdID(c)

			lr := cfg.LoadRepo(ifOpts.Workdir)
			show.AList("Current remotes", lr.RemoteLines())

			defRmt := cfg.DefaultRemote
			c.Infoln("Begin config remote URLs:")
			defUrl := interact.Ask("Please input URL for remote "+defRmt, "", nil)

			if len(defUrl) > 18 {
				err := lr.Cmd("remote", "set-url", defRmt, defUrl).Run()
				if err != nil {
					return err
				}
			} else {
				c.Infoln("input is invalid, skip init remote url.")
			}

			if !cfg.IsForkMode() {
				return lr.Cmd("remote", "-v").Run()
			}

			srcRmt := cfg.SourceRemote
			srcUrl := interact.Ask("Please input URL for remote "+srcRmt, "", nil)
			if len(srcUrl) > 18 {
				op := strutil.OrCond(lr.HasSourceRemote(), "set-url", "add")
				err = lr.Cmd("remote", op, srcRmt, srcUrl).Run()
				if err != nil {
					return err
				}
			} else {
				c.Infoln("input is invalid, skip init remote url.")
			}

			c.Infoln("\nNow, Git Remotes:")
			return lr.Cmd("remote", "-v").Run()
		},
	}
}
