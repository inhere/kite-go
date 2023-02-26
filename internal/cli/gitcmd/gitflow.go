package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/show"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/internal/biz/cmdbiz"
)

var CreatePRLink = &gcli.Command{
	Name:    "pr",
	Desc:    "create pull request link for current project",
	Aliases: []string{"pr-link"},
}

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
			ifOpts.BindCommonFlags1(c)
		},
		Func: func(c *gcli.Command, args []string) error {
			cfg := apputil.GitCfgByCmdID(c)

			lr := cfg.LoadRepo(ifOpts.Workdir)
			show.AList("Current remotes", lr.RemoteLines())

			c.Infof("config %q remote URL", cfg.DefaultRemote)
			defUrl := interact.Ask("URL", "", nil)

			if len(defUrl) > 18 {
				err := lr.Cmd("remote", "set-url", cfg.DefaultRemote, defUrl).Run()
				if err != nil {
					return err
				}
			}

			if !cfg.IsForkMode() {
				return lr.Cmd("remote", "-v").Run()
			}

			c.Infof("config %q remote URL", cfg.SourceRemote)
			srcUrl := interact.Ask("URL", "", nil)
			if len(defUrl) > 18 {
				err := lr.Cmd("remote", "set-url", cfg.SourceRemote, srcUrl).Run()
				if err != nil {
					return err
				}
			}

			return lr.Cmd("remote", "-v").Run()
		},
	}
}
