package gitx

import (
	"errors"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitwrap"
	"github.com/inherelab/kite/pkg/cmdutil"
	"github.com/inherelab/kite/pkg/gituse"
)

var (
	dryRun bool
	yesRun bool // Direct execution without confirmation

	interactive bool // interactively ask before executing command
)

// GitCommands commands for use git
var GitCommands = &gcli.Command{
	Name: "git",
	Desc: "tools for quick use `git` commands",
	Subs: []*gcli.Command{
		StatusInfo,
		RemoteInfo,
		AddCommitPush,
		AddCommitNotPush,
		TagCmd,
		gituse.OpenRemoteRepo,
		CreatePRLink,
		BatchPull,
	},
	Config: func(c *gcli.Command) {
		addListener(c)
	},
}

func addListener(c *gcli.Command) {
	c.On(gcli.EvtCmdOptParsed, func(obj ...interface{}) bool {
		c.Infoln("WorkDir:", c.WorkDir())
		return false
	})
	c.On(gcli.EvtCmdSubNotFound, func(data ...interface{}) (stop bool) {
		c.Errorln(data[1], "- the git subcommand is not exists, will call system command(TODO)")
		return true
	})
}

var StatusInfo = &gcli.Command{
	Name: "status",
	Aliases: []string{"st"},
	Desc: "git status command",
	Func: func(c *gcli.Command, args []string) error {
		return cmdutil.NewGitCmd("status").Run()
	},
}

var RemoteInfo = &gcli.Command{
	Name: "remote",
	Aliases: []string{"rmt"},
	Desc: "git remote command",
	Func: func(c *gcli.Command, args []string) error {
		err := gitwrap.New("remote", "-v").Run()
		if err != nil {
			return err
		}

		return nil
	},
}

var TagCmd = &gcli.Command{
	Name: "tag",
	Desc: "git tag commands",
	Subs: []*gcli.Command{
		TagCreate,
		TagDelete,
	},
}

var TagCreate = &gcli.Command{
	Name: "create",
	Aliases: []string{"new"},
	Desc: "create new tag by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}

var TagDelete = &gcli.Command{
	Name: "delete",
	Aliases: []string{"del", "rm", "remove"},
	Desc: "delete exists tags by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
