package gitcmd

import "github.com/gookit/gcli/v3"

// BranchCmd instance
var BranchCmd = &gcli.Command{
	Name:    "branch",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"br"},

	Subs: []*gcli.Command{
		DeleteBranch,
		CreateBranch,
		BranchListCmd,
	},
}

var BranchListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"ls"},
}

var CreateBranch = &gcli.Command{
	Name:    "new",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"n", "create"},
}

var DeleteBranch = &gcli.Command{
	Name:    "del",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"d", "rm", "delete"},
}
