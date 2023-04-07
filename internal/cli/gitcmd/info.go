package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var riOpts = struct {
	cmdbiz.CommonOpts
}{}

// RepoInfoCmd instance
var RepoInfoCmd = &gcli.Command{
	Name: "info",
	// Aliases: []string{"ls"},
	Desc: "show some info for the git repository",
	Config: func(c *gcli.Command) {
		riOpts.BindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		rp := gitw.NewRepo(riOpts.Workdir)

		show.AList("Information", rp.Info())
		return nil
	},
}

// RemoteInfoCmd instance
var RemoteInfoCmd = &gcli.Command{
	Name:    "remote",
	Aliases: []string{"rmt"},
	Desc:    "git remote command",
	Func: func(c *gcli.Command, args []string) error {
		err := gitw.New("remote", "-v").Run()
		if err != nil {
			return err
		}
		return nil
	},
}
