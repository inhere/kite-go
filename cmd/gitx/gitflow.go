package gitx

import "github.com/gookit/gcli/v3"

var GitFlow = &gcli.Command{
	Name: "gitflow",
	Desc: "tool commands for git-flow development",
	Aliases: []string{"gflow", "gf"},
	Subs: []*gcli.Command{
		InitFlow,
		CreatePRLink,
		BranchOperateEx,
	},
}

var CreatePRLink = &gcli.Command{
	Name: "pr-link",
	Desc: "create pull request link for current project",
	Aliases: []string{"pr"},
}

var InitFlow = &gcli.Command{
	Name: "init",
	Desc: "init repo remote and other info for current project",
}

var BranchOperateEx = &gcli.Command{
	Name: "branch",
	Desc: "checkout an new branch for development from `dist` remote",
	Aliases: []string{"br"},

	Subs: []*gcli.Command{
		DeleteBranch,
		CreateBranch,
	},
}

var CreateBranch = &gcli.Command{
	Name: "new",
	Desc: "checkout an new branch for development from `dist` remote",
	Aliases: []string{"n", "create"},
}

var DeleteBranch = &gcli.Command{
	Name: "del",
	Desc: "checkout an new branch for development from `dist` remote",
	Aliases: []string{"d", "rm", "delete"},
}
