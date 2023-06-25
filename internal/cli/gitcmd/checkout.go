package gitcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
)

var coOpts = struct {
	cmdbiz.CommonOpts
	NotPush bool `flag:"desc=not push to remote after checkout;shorts=n"`
}{}

// NewCheckoutCmd instance
func NewCheckoutCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "checkout",
		Desc:    "checkout to another branch and update to latest",
		Aliases: []string{"co", "switch"},
		Config: func(c *gcli.Command) {
			coOpts.BindCommonFlags(c)

			c.AddArg("branchName", "the target branch name", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			cfg := apputil.GitCfgByCmdID(c)
			rp := cfg.LoadRepo(coOpts.Workdir)

			branchName := c.Arg("branchName").String()

			defRemote := rp.DefaultRemote
			srcRemote := rp.SourceRemote
			defBranch := rp.DefaultBranch

			rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
				rr.DryRun = acpOpts.DryRun
				rr.Confirm = acpOpts.Confirm
				rr.OutToStd = true
			})

			// fetch remotes
			colorp.Infoln("Fetch remotes and check branch exists")
			rr.GitCmd("fetch", "--all", "-np")
			if err := rr.RunReset(); err != nil {
				return err
			}

			// local - checkout and pull
			if rp.HasLocalBranch(branchName) {
				colorp.Infof("Checkout branch %q from local and update it.\n", branchName)
				rr.GitCmd("checkout", branchName).
					GitCmd("pull", "-np", defRemote)

				if rp.HasSourceBranch(branchName) {
					rr.GitCmd("pull", srcRemote, branchName)
				}

				rr.GitCmd("pull", "-np", srcRemote, defBranch)

				if !coOpts.NotPush {
					if rp.UpstreamRemote() == defRemote {
						rr.GitCmd("push")
					} else {
						rr.GitCmd("push", "-u", defRemote, branchName)
					}
				}
				return rr.Run()
			}

			if rp.HasOriginBranch(branchName) {
				colorp.Infof("- checkout branch %q from remote %q\n", branchName, rp.DefaultRemote)
				// git checkout --track origin/NAME
				rr.GitCmd("checkout", "--track", rp.OriginBranch(branchName))

				if rp.HasSourceBranch(branchName) {
					rr.GitCmd("pull", srcRemote, branchName)
				}

				rr.GitCmd("pull", "-np", srcRemote, defBranch)
				return rr.Run()
			}

			if rp.HasSourceBranch(branchName) {
				colorp.Infof("- checkout branch %q from remote %q\n", branchName, rp.SourceRemote)
				rr.GitCmd("checkout", "-b", branchName, srcRemote+"/"+branchName).
					GitCmd("push", "-u", defRemote, branchName)
				return rr.Run()
			}

			return errorx.Rawf("want checkout branch %q not exists", branchName)
		},
	}
}
