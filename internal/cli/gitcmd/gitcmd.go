package gitcmd

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/gitx"
)

// GitOpts object
var GitOpts = cmdbiz.CommonOpts{}

// GitCommands commands for use git
var GitCommands = &gcli.Command{
	Name:    "git",
	Desc:    "tools for quick use `git` and more extra commands",
	Aliases: []string{"g", "gitx"},
	Subs: []*gcli.Command{
		RepoInfoCmd,
		// StatusInfoCmd,
		RemoteInfoCmd,
		NewCloneCmd(configProvider),
		NewAddCommitPush(configProvider),
		NewAddCommitCmd(configProvider),
		NewUpdateCmd(configProvider),
		NewUpdatePushCmd(configProvider),
		gitx.NewOpenRemoteCmd(nil),
		NewInitFlowCmd(),
		CreatePRLink,
		ShowLogCmd,
		ChangelogCmd,
		TagCmd,
		BatchCmd,
		NewBranchCmd(),
		NewGitEmojisCmd(),
	},
	Config: func(c *gcli.Command) {
		GitOpts.BindCommonFlags(c)

		c.On(events.OnCmdSubNotFound, SubCmdNotFound)
		c.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) bool {
			c.Infoln("[GIT] Workdir:", c.WorkDir())
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

	stop = true
	err := cmdr.NewGitCmd(name, args...).ToOSStdout().Run()
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
	cliutil.Warnf("%s: subcommand '%s' not found, will redirect to run `kite git %s`\n", pName, sName, sName)

	// dump.P(ctx.App.CommandNames(), ctx.App.HasCommand("git"))
	err := ctx.App.RunCmd("git", ctx.Cmd.RawArgs())
	if err != nil {
		dump.P(err)
	}

	return true
}
