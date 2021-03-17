package git

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/cmdutil"
)

var CmdsOfGit = &gcli.Command{
	Name: "git",
	Desc: "tools for quick use `git` commands",
	Subs: []*gcli.Command{
		StatusInfo,
		AddCommitPush,
		AddCommitNotPush,
	},
	Config: func(c *gcli.Command) {
		c.On(gcli.EvtCmdOptParsed, func(obj ...interface{}) {
			c.Infoln("workDir:", c.WorkDir())
		})
	},
}

var (
	dryRun bool

	acpOpts = struct {
		message string
		notPush bool
	}{}
)

var StatusInfo = &gcli.Command{
	Name: "status",
	Aliases: []string{"st"},
	Desc: "git status command",
	Func: func(c *gcli.Command, args []string) error {
		return cmdutil.NewGitCmd("status").Run()
	},
}

var AddCommitPush = &gcli.Command{
	Name: "acp",
	Desc: "run git add/commit/push at once command",
	Func: acpFunc,
	Config: func(c *gcli.Command) {
		AddCommitNotPush.Config(c)

		c.BoolVar(&acpOpts.notPush, &gcli.FlagMeta{
			Name:  "not-push",
			Alias: "np",
			Desc:  "dont execute git push",
		})
	},
}

var AddCommitNotPush = &gcli.Command{
	Name: "ac",
	Desc: "run git add/commit at once command",
	Func: acpFunc,
	Config: func(c *gcli.Command) {
		acpOpts.notPush = true // not push

		c.BoolOpt(&dryRun, "dry-run", "", false, "dont real execute command")
		c.StrOpt(&acpOpts.message, "message", "m", "", "the commit message")
		c.Required("message")

		c.BindArg(&gcli.Argument{
			Name:    "files",
			Desc:    "Only add special files. default will add all changed files",
			IsArray: true,
		})
	},
}

var acpFunc = func(c *gcli.Command, args []string) error {
	cr := cmdutil.NewRunner()
	cr.DryRun = dryRun

	if len(args) > 0 {
		cr.AddGitCmd("status", args...)
		cr.Addf("git add %s", strings.Join(args, " "))
	} else { // add all changed files
		cr.Add("git", "status")
		cr.AddLine("git add .")
	}

	cr.AddGitCmd("commit", "-m", acpOpts.message)

	if acpOpts.notPush == false {
		cr.AddGitCmd("push")
	}

	cr.Run()

	return nil
}
