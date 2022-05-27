package github

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/slog"
	"github.com/inherelab/kite/pkg/gitflow"
	"github.com/inherelab/kite/pkg/gituse"
)

// CmdForGithub commands
var CmdForGithub = &gcli.Command{
	Name:    "github",
	Aliases: []string{"gh", "gith", "hub", "ghub"},
	Desc:    "useful tools for use github",
	Subs: []*gcli.Command{
		gituse.OpenRemoteRepo,
		gitflow.UpdateCmd,
		gitflow.UpdatePushCmd,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdSubNotFound, func(data ...interface{}) (stop bool) {
			sub := data[1].(string)
			slog.Infof("subcommand '%s' not found in %s, redirect to git", sub, c.Name)

			c.App().RunCmd("git", c.RawArgs())
			return true
		})
	},
	// Hooks: map[string]gcli.HookFunc{
	// 	"": func(data ...interface{}) (stop bool) {
	// 		return false
	// 	},
	// },
}
