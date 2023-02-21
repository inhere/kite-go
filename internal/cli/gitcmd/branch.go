package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/cmdutil"
)

// BranchCmd instance
var BranchCmd = &gcli.Command{
	Name:    "branch",
	Desc:    "checkout an new branch for development from `source` remote",
	Aliases: []string{"br"},
	Subs: []*gcli.Command{
		BranchDeleteCmd,
		BranchCreateCmd,
		BranchListCmd,
		BranchSetupCmd,
	},
}

// NewBranchCmd instance
func NewBranchCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "branch",
		Desc:    "checkout an new branch for development from `source` remote",
		Aliases: []string{"br"},
		Subs: []*gcli.Command{
			BranchDeleteCmd,
			BranchCreateCmd,
			BranchListCmd,
			BranchSetupCmd,
		},
	}
}

// BranchListCmd instance
var BranchListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "checkout an new branch for development from `source` remote",
	Aliases: []string{"ls"},
	Func: func(c *gcli.Command, args []string) error {

		return errorx.New("TODO")
	},
}

var bcOpts = struct {
	cmdbiz.CommonOpts
	notToSrc bool
}{}

// BranchCreateCmd instance
var BranchCreateCmd = &gcli.Command{
	Name: "new",
	Desc: "create and checkout new branch for development",
	Help: `Workflow:
 1. git checkout DEFAULT_BRANCH
 2. git pull -np SOURCE_REMOTE DEFAULT_BRANCH
 3. git checkout -b NEW_BRANCH
 4. git push -u DEFAULT_REMOTE NEW_BRANCH
 5. git push SOURCE_REMOTE NEW_BRANCH`,
	Aliases: []string{"n", "create"},
	Config: func(c *gcli.Command) {
		bcOpts.BindCommonFlags(c)
		c.BoolOpt2(&bcOpts.notToSrc, "not-to-src, nts", "dont push new branch to the source remote")

		c.AddArg("branch", "the new branch name, allow vars: {ymd}", true)
	},
	Func: func(c *gcli.Command, args []string) error {
		cfg := apputil.GitCfgByCmdID(c)
		c.Infof("TIP: auto select git config type: %s(by cmd ID: %s)\n", cfg.HostType, c.ID())

		rp := cfg.LoadRepo(upOpts.Workdir)
		if err := rp.Check(); err != nil {
			return err
		}

		defRemote := rp.DefaultRemote
		srcRemote := rp.SourceRemote
		defBranch := rp.DefaultBranch
		newBranch := c.Arg("branch").String()

		rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
			rr.DryRun = acpOpts.DryRun
			rr.Confirm = acpOpts.Confirm
			rr.OutToStd = true
		})

		if rp.HasLocalBranch(newBranch) {
			c.Infof("TIP: local ")
			rr.GitCmd("checkout", newBranch).GitCmd("pull", "-np")
			return rr.Run()
		}

		curBranch := rp.CurBranchName()
		if defBranch != curBranch {
			rr.GitCmd("checkout", defBranch)
		}

		rr.GitCmd("pull", "-np", srcRemote, defBranch)
		rr.GitCmd("checkout", "-b", newBranch)
		rr.GitCmd("push", "-u", defRemote, newBranch)

		if !bcOpts.notToSrc {
			rr.GitCmd("push", srcRemote, newBranch)
		}
		return rr.Run()
	},
}

// BranchDeleteCmd instance
var BranchDeleteCmd = &gcli.Command{
	Name:    "del",
	Desc:    "checkout an new branch for development from `source` remote",
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
