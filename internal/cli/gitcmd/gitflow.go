package gitcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var GitFlow = &gcli.Command{
	Name: "gitflow",
	Desc: "tool commands for git-flow development",
	Subs: []*gcli.Command{
		InitFlow,
		CreatePRLink,
		UpdateNoPush,
		UpdateAndPush,
		BranchOperateEx,
	},

	Aliases: []string{"gflow", "gf"},
	Config: func(c *gcli.Command) {
		BindCommonOpts(c)
	},
}

var CreatePRLink = &gcli.Command{
	Name:    "pr-link",
	Desc:    "create pull request link for current project",
	Aliases: []string{"pr"},
}

var InitFlow = &gcli.Command{
	Name: "init",
	Desc: "init repo remote and other info for current project",
}

var (
	gfUpOpts = struct {
		push bool
	}{}
	UpdateNoPush = &gcli.Command{
		Name: "update",
		Desc: "Update codes from origin and main remote repositories",
		Func: handleUpdatePush,

		Aliases: []string{"up"},
		Config: func(c *gcli.Command) {
			c.BoolOpt(&gfUpOpts.push, "push", "p", false, "Push to origin remote after update")
		},
	}

	UpdateAndPush = &gcli.Command{
		Name: "update-push",
		Desc: "Update codes from origin and main remote repositories, then push to remote",
		Func: handleUpdatePush,

		Aliases: []string{"upp", "up-push"},
	}
)

func handleUpdatePush(c *gcli.Command, args []string) error {
	return errors.New("TODO")
}

var BranchOperateEx = &gcli.Command{
	Name:    "branch",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"br"},

	Subs: []*gcli.Command{
		DeleteBranch,
		CreateBranch,
	},
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
