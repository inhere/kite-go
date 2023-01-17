package gitx

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite/pkg/gituse"
)

var (
	dryRun  bool
	yesRun  bool // Direct execution without confirmation
	workdir string

	confirm bool // interactively ask before executing command
)

func bindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&dryRun, "dry-run", "", false, "run workflow, but dont real execute command")
	c.StrOpt(&workdir, "workdir", "w", "", "the command workdir path")
}

// GitCommands commands for use git
var GitCommands = &gcli.Command{
	Name:    "git",
	Desc:    "tools for quick use `git` commands",
	Aliases: []string{"g"},
	Subs: []*gcli.Command{
		StatusInfo,
		RemoteInfo,
		AddCommitPush,
		AddCommitNotPush,
		TagCmd,
		UpdateCmd,
		gituse.OpenRemoteRepo,
		CreatePRLink,
		BatchCmd,
		Changelog,
		ShowLog,
		InitFlow,
		UpdateNoPush,
		UpdateAndPush,
		BranchOperateEx,
	},
	Config: func(c *gcli.Command) {
		addListener(c)

		c.BoolOpt(&dryRun, "dry-run", "", false, "Dry-run the workflow, dont real execute")
		c.BoolOpt(&yesRun, "yes", "y", false, "Direct execution without confirmation")
	},
}

func addListener(c *gcli.Command) {
	c.On(gcli.EvtCmdOptParsed, func(ctx *gcli.HookCtx) bool {
		c.Infoln("Workdir:", c.WorkDir())
		return false
	})
	c.On(gcli.EvtCmdSubNotFound, func(ctx *gcli.HookCtx) (stop bool) {
		c.Errorln(ctx.Str("name"), "- the git subcommand is not exists, will call system command(TODO)")
		return true
	})
}

var StatusInfo = &gcli.Command{
	Name:    "status",
	Aliases: []string{"st"},
	Desc:    "git status command",
	Func: func(c *gcli.Command, args []string) error {
		return cmdr.NewGitCmd("status").Run()
	},
}

var RemoteInfo = &gcli.Command{
	Name:    "remote",
	Aliases: []string{"rmt"},
	Desc:    "git remote command",
	Func: func(c *gcli.Command, args []string) error {
		err := gitw.New("remote", "-v").Run()
		if err != nil {
			return err
		}

		return nil
	},
}
