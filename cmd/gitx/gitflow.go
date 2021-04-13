package gitx

import "github.com/gookit/gcli/v3"

var GitFlow = &gcli.Command{
	Name: "gitflow",
	Desc: "tool commands for git-flow development",
	Aliases: []string{"gf"},
	Subs: []*gcli.Command{
		CreatePRLink,
		CreateBranch,
	},
}

var CreatePRLink = &gcli.Command{
	Name: "pr-link",
	Desc: "create pull request link for current project",
	Aliases: []string{"pr"},
}

var CreateBranch = &gcli.Command{
	Name: "new-branch",
	Desc: "checkout an new branch for development from `dist` remote",
	Aliases: []string{"nb", "cobr"},
}
