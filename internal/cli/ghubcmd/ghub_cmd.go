package ghubcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/slog"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/gitflow"
)

// CmdForGithub commands
var CmdForGithub = &gcli.Command{
	Name:    "github",
	Aliases: []string{"gh", "gith", "hub", "ghub"},
	Desc:    "useful tools for use github",
	Subs: []*gcli.Command{
		gitx.NewOpenRemoteCmd(gitx.GithubHost),
		gitflow.UpdateCmd,
		gitflow.UpdatePushCmd,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdSubNotFound, func(ctx *gcli.HookCtx) (stop bool) {
			sub := ctx.Str("name")
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
