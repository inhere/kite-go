package gitcmd

import (
	"github.com/gookit/gcli/v3"
)

var CreatePRLink = &gcli.Command{
	Name:    "pr",
	Desc:    "create pull request link for current project",
	Aliases: []string{"pr-link"},
}

var InitFlow = &gcli.Command{
	Name: "init",
	Desc: "init repo remote and other info for current project",
}
