package gitcmd

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/gitx"
)

// GitOpts object
var GitOpts = struct {
	AutoGitDir
	cmdbiz.CommonOpts
}{}

// GitCommands commands for use git
var GitCommands = &gcli.Command{
	Name:    "git",
	Desc:    "tools for quick use `git` and more extra commands",
	Aliases: []string{"g", "gitx"},
	Subs: []*gcli.Command{
		RepoInfoCmd,
		// StatusInfoCmd,
		RemoteInfoCmd,
		NewPullRequestCmd(),
		NewCloneCmd(configProvider),
		NewAddCommitPush(),
		NewAddCommitCmd(),
		NewUpdateCmd(),
		NewUpdatePushCmd(),
		NewOpenRemoteCmd(configProvider),
		NewInitFlowCmd(),
		NewCheckoutCmd(),
		ShowLogCmd,
		ChangelogCmd,
		TagCmd,
		BatchCmd,
		NewBranchCmd(),
		NewGitEmojisCmd(),
	},
	Config: func(c *gcli.Command) {
		GitOpts.BindCommonFlags(c)
		GitOpts.BindChdirFlags(c)

		c.On(events.OnCmdSubNotFound, SubCmdNotFound)
		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) bool {
			c.Infoln("[kite.GIT] Workdir:", c.WorkDir())
			return false
		})
	},
}

func configProvider() *gitx.Config {
	return app.Gitx()
}

// SubCmdNotFound handle
func SubCmdNotFound(ctx *gcli.HookCtx) (stop bool) {
	pName := ctx.Cmd.Name
	name := ctx.Str("name")
	args := ctx.Strings("args")

	if name[0] == '@' {
		rawNa := name
		name = name[1:]
		cliutil.Warnf("%s: %s - start with @, will be direct call `git %s` on system\n", pName, rawNa, name)
	} else {
		cliutil.Warnf("%s: subcommand %q is not found, will call `git %s` on system\n", pName, name, name)
	}

	colorp.Infoln("[kite.GIT] Workdir:", ctx.Cmd.WorkDir())

	stop = true
	err := cmdr.NewGitCmd(name, args...).PrintCmdline().ToOSStdout().Run()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if ok && ee.ProcessState != nil {
			st := ee.Sys().(syscall.WaitStatus)
			if st.Signaled() {
				fmt.Println("Quited.")
				return
			}
		}

		cliutil.Errorln("Exec Error:", err)
	}
	return
}

// RedirectToGitx handle
func RedirectToGitx(ctx *gcli.HookCtx) (stop bool) {
	if ctx.App == nil {
		return
	}

	pName := ctx.Cmd.Name
	sName := ctx.Str("name")
	args := ctx.Cmd.RawArgs()
	cliutil.Warnf("%s: subcommand '%s' not found, redirect to run `kite git %s`, args: %v\n", pName, sName, sName, args[1:])
	cmdbiz.ProxyCC.AutoSetByName(pName, sName, args[1:])

	// dump.P(ctx.App.CommandNames(), ctx.App.HasCommand("git"))
	err := ctx.App.RunCmd("git", args)
	if err != nil {
		colorp.Errorln(err)
	}
	return true
}

// AutoGitDir auto change to .git dir
type AutoGitDir struct {
	// AutoGit auto find .git dir in parent.
	AutoGit bool
	GitHost string
}

// BindChdirFlags for auto change dir
func (a *AutoGitDir) BindChdirFlags(c *gcli.Command) {
	wd := c.WorkDir()

	c.BoolOpt2(&a.AutoGit, "auto-root, auto-git", "auto find .git dir in parent and chdir to it", gflag.WithValidator(func(val string) error {
		if strutil.QuietBool(val) {
			repoDir, changed := fsutil.SearchNameUpx(wd, gitw.GitDir)
			if changed {
				goutil.MustOK(c.ChWorkDir(repoDir))
				cliutil.Yellowf("NOTICE: auto founded git root and will chdir to: %s\n", repoDir)
			}
		}
		return nil
	}))
}
