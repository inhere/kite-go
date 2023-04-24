package gitcmd

import (
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/basefn"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
	"github.com/inhere/kite-go/pkg/gitx"
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

var blOpts = struct {
	cmdbiz.CommonOpts
	Remote string `flag:"desc=only show branches on the remote;shorts=r"`
	Match  string `flag:"desc=the branch name match pattern;shorts=p,m"`
	Regex  bool   `flag:"desc=enable regex for match pattern;shorts=reg"`
	All    bool   `flag:"desc=display all branches;shorts=a"`
	Exec   string `flag:"desc=execute command for each branch;shorts=x"`
	Delete bool   `flag:"desc=delete matched branches;shorts=d"`
}{}

// BranchListCmd instance
var BranchListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "list or search branches on local or remote",
	Aliases: []string{"search", "ls"},
	Examples: `
# list branches by glob pattern
{$fullCmd} -m "fea*"

# list branches by regex pattern
{$fullCmd} --reg -m "fea_\d+"

# find and delete remote branches
{$fullCmd} -r origin -x "git push {remote} --delete {branch}"
`,
	Config: func(c *gcli.Command) {
		blOpts.BindCommonFlags(c)
		c.MustFromStruct(&blOpts, gflag.TagRuleNamed)
		c.AddArg("match", "the branch name match pattern, same as --match|-m|-p")
	},
	Func: func(c *gcli.Command, args []string) error {
		rp := app.Gitx().LoadRepo(blOpts.Workdir)
		if err := rp.FetchAll("-np"); err != nil {
			return err
		}

		colorp.Infoln("Load repo branches ...")
		bis := rp.BranchInfos()

		tle := "Local"
		var brs []*gitw.BranchInfo

		if blOpts.All {
			tle = "Local+Remotes"
			brs = bis.All()
		} else if blOpts.Remote != "" {
			tle = blOpts.Remote
			brs = bis.Remotes(blOpts.Remote)
		} else {
			brs = bis.Locales()
		}

		var number int
		matcher := gitx.NewBranchMatcher(blOpts.Match, blOpts.Regex)

		colorp.Infof("Branches on %s\n", tle)
		for _, info := range brs {
			if matcher.Match(info.Short) {
				number++
				colorp.Infof(" - %16s %s\n", info.Name, info.HashMsg)
			}
		}

		colorp.Infoln("Match Total:", number)
		return nil
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

var bdOpts = struct {
	cmdbiz.CommonOpts
	Remote string `flag:"desc=remote name for delete branches;shorts=r"`
}{}

// BranchDeleteCmd instance
var BranchDeleteCmd = &gcli.Command{
	Name:    "del",
	Desc:    "delete one or more branches from local and remote",
	Aliases: []string{"d", "rm", "delete"},
	Func: func(c *gcli.Command, args []string) error {
		// git branch -d feature-*
		// git push origin --delete feature-*
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
