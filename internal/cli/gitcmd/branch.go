package gitcmd

import (
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/basefn"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
)

// NewBranchCmd instance
func NewBranchCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "branch",
		Desc:    "git branch commands extension",
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

		rp := cfg.LoadRepo(upOpts.Workdir)
		if err := rp.Check(); err != nil {
			return err
		}

		defRemote := rp.DefaultRemote
		srcRemote := rp.SourceRemote
		defBranch := rp.DefaultBranch
		newBranch := c.Arg("branch").String()
		if strings.Contains(newBranch, "{ymd}") {
			newBranch = strings.Replace(newBranch, "{ymd}", timex.Now().DateFormat("ymd"), -1)
		}

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

		// fetch remotes and check branch exists
		colorp.Infoln("Fetch remotes and check branch exists")
		rr.GitCmd("fetch", "-np", defRemote).GitCmd("fetch", "-np", srcRemote)
		if err := rr.Run(); err != nil {
			return err
		}

		rr.Reset()

		if rp.HasOriginBranch(newBranch) {
			colorp.Warnf("the branch %q already exists in remote %q\n", newBranch, defRemote)
			return rp.QuickRun("checkout", newBranch)
		}

		if rp.HasSourceBranch(newBranch) {
			colorp.Warnf("the branch %q already exists in remote %q\n", newBranch, srcRemote)
			rr.GitCmd("checkout", newBranch).GitCmd("push", "-u", defRemote, newBranch)
			return rr.Run()
		}

		// do checkout new branch and push to remote
		colorp.Infoln("Do checkout new branch and push to remote")
		curBranch := rp.CurBranchName()
		if defBranch != curBranch {
			rr.GitCmd("checkout", defBranch)
		}

		rr.GitCmd("pull", srcRemote, defBranch)
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

var bsOpts = struct {
	cmdbiz.CommonOpts
	notToSrc bool
}{}

// BranchSetupCmd instance
var BranchSetupCmd = &gcli.Command{
	Name: "setup",
	Desc: "setup a new checkout branch on fork develop mode",
	Help: `Workflow:
git fetch DEFAULT_REMOTE

if DEFAULT_REMOTE/BRANCH exist:
	git branch --set-upstream-to=DEFAULT_REMOTE/BRANCH
else:
	git push --set-upstream DEFAULT_REMOTE BRANCH

git push SOURCE_REMOTE BRANCH
`,
	Aliases: []string{"init"},
	Config: func(c *gcli.Command) {
		bsOpts.BindCommonFlags(c)
		c.BoolOpt2(&bsOpts.notToSrc, "not-to-src, nts", "dont push branch to the source remote")
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		cfg := apputil.GitCfgByCmdID(c)
		rp := cfg.LoadRepo(bsOpts.Workdir)
		if err := rp.Check(); err != nil {
			return err
		}

		if err := rp.FetchOrigin(); err != nil {
			return err
		}

		// git fetch origin
		// exist: git branch --set-upstream-to=DEFAULT_REMOTE/BRANCH
		// not exist: git push --set-upstream DEFAULT_REMOTE BRANCH
		brName := rp.CurBranchName()
		if rp.HasOriginBranch(brName) {
			err = rp.Cmd("branch").Argf("--set-upstream-to=%s/%s", rp.DefaultRemote, brName).Run()
		} else {
			err = rp.Cmd("push", "-u", rp.DefaultRemote, brName).Run()
		}

		return basefn.CallOn(!bsOpts.notToSrc, func() error {
			return rp.Cmd("push", rp.SourceRemote, brName).Run()
		})
	},
}
