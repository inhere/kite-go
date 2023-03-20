package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
	"github.com/inhere/kite-go/pkg/gitx"
)

var upOpts = struct {
	cmdbiz.CommonOpts
	notPush bool
}{}

const upHelp = `
Workflow:
 1. pre-check for update
 2. git pull default_remote
On fork_mode=true:
 3. git pull source_remote current_branch
 4. git pull source_remote default_branch
`

// NewUpdatePushCmd instance
func NewUpdatePushCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	return &gcli.Command{
		Name:    "update-push",
		Desc:    "Update code from remotes, then push to default remote",
		Help:    upHelp,
		Aliases: []string{"up-push", "upp"},
		Config: func(c *gcli.Command) {
			upOpts.BindCommonFlags(c)

			c.BoolVar(&upOpts.notPush, &gcli.FlagMeta{
				Name:   "not-push",
				Desc:   "dont push update to default remote",
				Shorts: []string{"np"},
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			return updateHandleFunc(c, args, cfgGetter())
		},
	}
}

// NewUpdateCmd instance
func NewUpdateCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	return &gcli.Command{
		Name:    "update",
		Desc:    "Update code from remote repositories",
		Help:    upHelp,
		Aliases: []string{"up"},
		Config: func(c *gcli.Command) {
			upOpts.BindCommonFlags(c)
		},
		Func: func(c *gcli.Command, args []string) error {
			upOpts.notPush = true
			return updateHandleFunc(c, args, cfgGetter())
		},
	}
}

func updateHandleFunc(c *gcli.Command, _ []string, cfg *gitx.Config) (err error) {
	rp := cfg.LoadRepo(upOpts.Workdir)
	c.Infoln("TIP: pre-check for update repository data")

	defRemote := rp.DefaultRemote
	srcRemote := rp.SourceRemote
	curBranch := rp.CurBranchName()

	if !rp.HasDefaultRemote() {
		return c.NewErrf(
			"not found default remote %q, please add it by `git remote add %s URL`",
			defRemote, defRemote)
	}

	defBranch := rp.DefaultBranch
	upstream := rp.UpstreamRemote()

	if !rp.IsDefaultRemote(upstream) {
		c.Warnf("TIP: current upstream remote is not %q, will auto update it.\n", defRemote)
		if !rp.HasRemoteBranch(curBranch, defRemote) {
			err = rp.Cmd("push", "-u", defRemote).Run()
		} else {
			err = rp.SetUpstreamTo(defRemote, curBranch)
		}

		if err != nil {
			return err
		}
	}

	if !rp.IsForkMode() {
		c.Infoln("Do update repository data from remote ...")
		return rp.Cmd("pull", "-np").Run()
	}

	rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
		rr.DryRun = upOpts.DryRun
		// rr.Confirm = confirm
		rr.OutToStd = true
	})

	// update from default remote
	rr.GitCmd("pull", "-np")

	if !rp.HasSourceRemote() {
		c.Warnf(
			"TIP: the source remote %q is not added, please add by `git remote add %s URL`\n",
			srcRemote, srcRemote,
		)
	} else {
		// update from source remote current_branch
		if curBranch != defBranch && rp.HasRemoteBranch(curBranch, srcRemote) {
			rr.GitCmd("pull", srcRemote, curBranch)
		}

		// pull latest from source remote DefaultBranch
		rr.GitCmd("pull", "-np", srcRemote, defBranch)
	}

	// push to default remote
	if !upOpts.notPush {
		rr.GitCmd("push")
	}

	c.Infoln("Do update repository data from remotes ...")
	return rr.Run()
}
