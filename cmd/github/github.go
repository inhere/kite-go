package github

import (
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/gituse"
)

var CmdForGithub = &gcli.Command{
	Name: "github",
	Aliases: []string{"gh", "hub", "ghub"},
	Desc: "useful tools for use github",
	Subs: []*gcli.Command{
		gituse.OpenRemoteRepo,
	},
}
