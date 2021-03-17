package git

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/pkg/cmdutil"
)

var acpOpts = struct {
	message string
	notPush bool
}{}

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
	runPush := acpOpts.notPush == false
	if c.Name == "ac" {
		runPush = false // not push
	}

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

	if runPush {
		cr.AddGitCmd("push")
	}

	cr.Run()
	return nil
}
