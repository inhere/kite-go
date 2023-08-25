package gitcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
)

var upOpts = struct {
	cmdbiz.CommonOpts
	notPush  bool
	fetchAll bool
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
func NewUpdatePushCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "update-push",
		Desc:    "Update code from remotes, then push to default remote",
		Help:    upHelp,
		Aliases: []string{"up-push", "upp"},
		Config: func(c *gcli.Command) {
			upOpts.BindCommonFlags(c)

			c.BoolVar(&upOpts.notPush, &gcli.CliOpt{
				Name:   "not-push",
				Desc:   "dont push update to default remote",
				Shorts: []string{"np"},
			})
			c.BoolOpt2(&upOpts.fetchAll, "fetch-all,fetch,f", "only fetch all remotes with option -np, dont run pull")
		},
		Func: func(c *gcli.Command, args []string) error {
			return updateHandleFunc(c, args)
		},
	}
}

// NewUpdateCmd instance
func NewUpdateCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "update",
		Desc:    "Update code from remote repositories",
		Help:    upHelp,
		Aliases: []string{"up"},
		Config: func(c *gcli.Command) {
			upOpts.BindCommonFlags(c)
			c.BoolOpt2(&upOpts.fetchAll, "fetch-all,fetch,fa", "only fetch all remotes with option -np, dont run pull")
		},
		Func: func(c *gcli.Command, args []string) error {
			upOpts.notPush = true
			return updateHandleFunc(c, args)
		},
	}
}

func updateHandleFunc(c *gcli.Command, _ []string) (err error) {
	cfg := apputil.GitCfgByCmdID(c)
	rp := cfg.LoadRepo(upOpts.Workdir)
	rp.SetDryRun(upOpts.DryRun)
	c.Infoln("TIP: pre-check for update repository data")

	// update remote info
	err = rp.Cmd("fetch", "--all", "-np").Run()
	if err != nil {
		return err
	}
	if upOpts.fetchAll {
		return nil
	}

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
		colorp.Infoln("Do update repository data from remote ...")
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
		colorp.Warnf(
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

	colorp.Infoln("Do update repository data from remotes ...")
	return rr.Run()
}

func updateBranch() error {
	// TODO ...
	return nil
}
