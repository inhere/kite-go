package gitcmd

import (
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/brinfo"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/basefn"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite-go/internal/app"
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

var bdOpts = struct {
	cmdbiz.CommonOpts
	Fetch  bool   `flag:"desc=fetch remote branches before search;shorts=f"`
	Remote string `flag:"desc=only find branches on the remote;shorts=r"`
}{}

// BranchDeleteCmd quickly delete branches
var BranchDeleteCmd = &gcli.Command{
	Name:    "delete",
	Desc:    "delete branches",
	Aliases: []string{"del", "rm"},
	Examples: `
# Delete local branches
{$fullCmd} fix_1 fix_2 fea_1
{$fullCmd} "fea*" # use glob match
{$fullCmd} "prefix:fea" # use prefix match
{$fullCmd} "suffix:_dev" # use suffix match
{$fullCmd} "regex:fix_\d+" # use regex match

# Delete remote branches
{$fullCmd} -r origin fix_* # use glob match
`,
	Config: func(c *gcli.Command) {
		bdOpts.BindCommonFlags(c)
		c.MustFromStruct(&bdOpts)
		c.AddArg("branches", "the branch names or match pattern(can with match mode)", true, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		rp := app.Gitx().LoadRepo(blOpts.Workdir)
		if blOpts.Fetch {
			if err := rp.FetchAll("-np"); err != nil {
				return err
			}
		}

		matcher := brinfo.NewMulti()
		patterns := c.Arg("branches").Strings()

		var branches []string
		for _, pattern := range patterns {
			if gitutil.IsBranchName(pattern) {
				branches = append(branches, pattern)
				continue
			}

			// as match pattern
			matcher.Add(brinfo.NewMatcher(pattern))
		}

		bis := rp.BranchInfos()
		opt := &gitw.SearchOpt{Flag: gitw.BrSearchLocal, Remote: bdOpts.Remote}
		if bdOpts.Remote != "" {
			opt.Flag = gitw.BrSearchRemote
		}

		brs := bis.SearchV2(matcher, opt)
		num := len(brs)

		rp.SetDryRun(blOpts.DryRun)
		colorp.Cyanf("Will delete %d branches, Glob match number %d\n", num+len(branches), num)

		if blOpts.Confirm && !cliutil.Confirm("Do you want to continue?") {
			colorp.Cyanln("Canceled")
			return nil
		}

		for _, info := range brs {
			// git branch -d feature-*
			// git push origin --delete feature-*
			colorp.Warnf("Do delete branch: %s\n", info.Name)
			if err := rp.BranchDelete(info.Short, info.Remote); err != nil {
				return err
			}
		}
		return nil
	},
}

var blOpts = struct {
	cmdbiz.CommonOpts
	Remote string `flag:"desc=only show branches on the remote;shorts=r"`
	Match  string `flag:"desc=the branch name match pattern, default use glob mode;shorts=p,m"`
	Limit  int    `flag:"desc=limit the max match branches;shorts=n"`
	All    bool   `flag:"desc=display all branches;shorts=a"`
	Exec   string `flag:"desc=execute command for each branch;shorts=x"`
	Delete bool   `flag:"desc=delete matched branches;shorts=d"`
	Fetch  bool   `flag:"desc=fetch remote branches before search;shorts=f"`
	Update bool   `flag:"desc=update branches after exec/delete;shorts=u"`
}{}

// BranchListCmd instance
var BranchListCmd = &gcli.Command{
	Name:    "list",
	Desc:    "list or search branches on local or remote",
	Aliases: []string{"search", "ls"},
	Examples: `
# List branches by glob pattern
{$fullCmd} -m "fea*"
{$fullCmd} -m "glob:fix_2[1|2]*"
# List branches by prefix pattern
{$fullCmd} -m "prefix:fea-"

# List branches by regex pattern
{$fullCmd} -m "regex:fea_\d+"
# match like fix_23_04
{$fullCmd} -m "regex:f[a-z]*_[0-9]+_\d+"

# Find and delete remote branches
{$fullCmd} -r origin -m fix_* --delete
# or
{$fullCmd} -r origin -m fix_* -x "git push {remote} --delete {branch}"
`,
	Config: func(c *gcli.Command) {
		blOpts.BindCommonFlags(c)
		c.MustFromStruct(&blOpts, gflag.TagRuleNamed)
		c.AddArg("match", "the branch name match pattern, same as --match|-m|-p").WithAfterFn(func(arg *gcli.CliArg) error {
			blOpts.Match = arg.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, args []string) error {
		rp := app.Gitx().LoadRepo(blOpts.Workdir)

		if blOpts.Fetch {
			if err := rp.FetchAll("-np"); err != nil {
				return err
			}
		}

		colorp.Infoln("Load repo branches ...")
		bis := rp.BranchInfos()

		if blOpts.Match == "" {
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

			colorp.Cyanf("Branches on %q(total: %d)\n", tle, len(brs))
			for i, info := range brs {
				colorp.Infof(" %-24s %s\n", info.Short, info.HashMsg)
				if blOpts.Limit > 0 && i+1 >= blOpts.Limit {
					break
				}
			}
			return nil
		}

		// ------ search branches

		tle := "Local"
		if blOpts.All {
			tle = "Local+Remotes"
		} else if blOpts.Remote != "" {
			tle = blOpts.Remote
		}

		matcher := brinfo.NewMatcher(blOpts.Match)

		opt := &gitw.SearchOpt{Flag: gitw.BrSearchLocal, Remote: blOpts.Remote, Limit: blOpts.Limit}
		brs := bis.SearchV2(matcher, opt)
		colorp.Cyanf("Branches on %q(found: %d, %s)\n", tle, len(brs), matcher.String())

		if blOpts.Confirm && !cliutil.Confirm("Do you want to continue?") {
			colorp.Cyanln("Canceled")
			return nil
		}

		spl := textutil.NewVarReplacer("{,}")
		rp.SetDryRun(blOpts.DryRun)

		for _, info := range brs {
			if blOpts.Delete {
				// git branch -d feature-*
				// git push origin --delete feature-*
				colorp.Warnf("Do delete branch: %s\n", info.Name)
				if err := rp.BranchDelete(info.Short, info.Remote); err != nil {
					return err
				}
			} else if blOpts.Exec != "" {
				vs := map[string]string{
					"branch": info.Short,
					"hash":   info.Hash,
					"remote": info.Remote,
				}

				execCmd := cmdr.NewCmdline(spl.RenderSimple(blOpts.Exec, vs)).
					WithWorkDir(rp.Dir()).
					WithDryRun(blOpts.DryRun).
					PrintCmdline()

				if err := execCmd.Run(); err != nil {
					return err
				}
			} else {
				colorp.Infof(" %-24s %s\n", info.Short, info.HashMsg)
			}
		}

		if blOpts.Update {
			if err := rp.FetchAll("-np"); err != nil {
				return err
			}
		}
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
		rr.GitCmd("fetch", "--all", "-np")
		if err := rr.RunReset(); err != nil {
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
