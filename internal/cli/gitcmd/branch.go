package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite/pkg/gitx"
)

// BranchCmd instance
var BranchCmd = &gcli.Command{
	Name:    "branch",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"br"},

	Subs: []*gcli.Command{
		BranchDeleteCmd,
		BranchCreateCmd,
		BranchListCmd,
		BranchSetupCmd,
	},
}

// BranchListCmd instance
var BranchListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"ls"},
	Func: func(c *gcli.Command, args []string) error {

		return errorx.New("TODO")
	},
}

var bcOpts = struct {
	gitx.CommonOpts
	notToSrc bool
}{}

// BranchCreateCmd instance
var BranchCreateCmd = &gcli.Command{
	Name: "new",
	Desc: "checkout an new branch for development from `dist` remote",
	Help: `Workflow:
 1. git checkout to DEFAULT_BRANCH
 2. git pull main DEFAULT_BRANCH
 3. git checkout -b NEW_BRANCH
 4. git push -u DEFAULT_REMOTE NEW_BRANCH
 5. git push SOURCE_REMOTE NEW_BRANCH`,
	Aliases: []string{"n", "create"},
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&bcOpts.notToSrc, "not-to-src, nts", "dont push new branch to the source remote")
	},
	Func: func(c *gcli.Command, args []string) error {
		return errorx.New("TODO")
	},
}

// BranchDeleteCmd instance
var BranchDeleteCmd = &gcli.Command{
	Name:    "del",
	Desc:    "checkout an new branch for development from `dist` remote",
	Aliases: []string{"d", "rm", "delete"},
	Func: func(c *gcli.Command, args []string) error {

		return errorx.New("TODO")
	},
}

// BranchSetupCmd instance
var BranchSetupCmd = &gcli.Command{
	Name:    "setup",
	Desc:    "setup a new checkout branch on fork develop mode",
	Help:    `Will setup upstream remote to default_remote`,
	Aliases: []string{"init"},
	Func: func(c *gcli.Command, args []string) error {
		// git fetch origin
		// exist: git br --set-upstream-to=origin/br
		// not exist: git push -u origin
		return errorx.New("TODO")
	},
}
