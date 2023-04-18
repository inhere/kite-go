package glabcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/cmdutil"
)

var rcOpts cmdbiz.CommonOpts

// ResolveConflictCmd instance
var ResolveConflictCmd = &gcli.Command{
	Name: "resolve",
	Desc: "Resolve conflicts preparing for current git branch.",
	Help: `Workflow:
1. will checkout to <cyan>branch</cyan>
2. will update code by <cyan>git pull</cyan>
3. update the <cyan>branch</cyan> codes from source repository
4. merge current-branch codes from source repository
5. please resolve conflicts by tools or manual
`,
	Aliases: []string{"rc"},
	Config: func(c *gcli.Command) {
		rcOpts.BindCommonFlags(c)

		c.AddArg("branch", "The conflicts target branch name. eg: qa, pre, master", true)
	},
	Func: func(c *gcli.Command, args []string) error {
		gl := app.Glab()
		lr := gl.LoadRepo(rcOpts.Workdir)

		curBr := lr.CurBranchName()
		br := c.Arg("branch").String()
		br = lr.ResolveBranch(br)

		rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
			rr.Workdir = rcOpts.Workdir
			rr.DryRun = rcOpts.DryRun
		})

		// fetch all remotes to latest
		rr.GitCmd("fetch", "-np", gl.DefaultRemote)
		rr.GitCmd("fetch", "-np", gl.SourceRemote)

		if lr.HasLocalBranch(br) {
			rr.GitCmd("checkout", br)

			if lr.HasOriginBranch(br) {
				rr.GitCmd("pull", gl.DefaultRemote, br)
			} else {
				rr.GitCmd("push", "-u", gl.DefaultRemote, br)
			}

			// update the branch codes from source remote
			rr.GitCmd("pull", gl.SourceRemote, br)
		} else if lr.HasOriginBranch(br) {
			// git checkout --track origin/NAME
			rr.GitCmd("checkout", "--track", gl.OriginBranch(br)).
				GitCmd("pull", gl.SourceRemote, br)
		} else if lr.HasSourceBranch(br) {
			rr.GitCmd("checkout", "--track", gl.SourceBranch(br)).
				GitCmd("push", "-u", gl.DefaultRemote)
		} else {
			return c.NewErrf("branch %q not found in local and remotes", br)
		}

		if err := rr.Run(); err != nil {
			return err
		}
		c.Println("---------------------------------------------------------------------")

		// pull target branch from source remote, will merge to current branch
		_ = lr.Cmd("pull", gl.SourceRemote, curBr).Run()

		c.Println("---------------------------------------------------------------------")
		colorp.Successln("Now, please resolve conflicts by tools or manual ...")
		colorp.Infoln("TIP 1: If you want stop resolve, can run:\n  git merge --abort && git checkout", curBr)
		colorp.Warnln("TIP 2: Can execute follow command after resolved for quick commit:")
		colorp.Infoln("  git add . && git commit && git push && kite gl pr -o head && git checkout", curBr)
		return nil
	},
}
