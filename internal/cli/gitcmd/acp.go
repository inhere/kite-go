package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/inhere/kite/pkg/cmdutil"
	"github.com/inhere/kite/pkg/gitx"
)

var acpOpts = struct {
	gitx.CommonOpts
	message string
	notPush bool
}{}

const acpHelp = `  Commit types:
 build     "Build system"
 chore     "Chore"
 ci        "CI"
 docs      "Documentation"
 feat      "Features"
 fix       "Bug fixes"
 perf      "Performance"
 refactor  "Refactor"
 style     "Style"
 test      "Testing"
`

// AddCommitPush command
var AddCommitPush = &gcli.Command{
	Name: "acp",
	Desc: "run `git add/commit/push` at once command",
	Func: acpHandleFunc,
	Help: acpHelp,
	Config: func(c *gcli.Command) {
		AddCommitNotPush.Config(c)

		c.BoolOpt(&acpOpts.notPush, "not-push", "np", false, "dont execute git push")
	},
}

// AddCommitNotPush command
var AddCommitNotPush = &gcli.Command{
	Name: "ac",
	Desc: "run git add/commit at once command",
	Help: acpHelp,
	Func: acpHandleFunc,
	Config: func(c *gcli.Command) {
		acpOpts.BindCommonFlags(c)

		c.BoolOpt(&confirm, "interactive", "i", false, "confirm ask before executing command")
		c.StrOpt2(&acpOpts.message, "message,m", "the git commit message", gflag.WithRequired())

		c.BindArg(&gcli.CliArg{
			Name:    "files",
			Desc:    "Only add special files. default will add all changed files",
			Arrayed: true,
		})
	},
}

func acpHandleFunc(c *gcli.Command, args []string) error {
	runPush := acpOpts.notPush == false
	if c.Name == "ac" {
		runPush = false // not push
	}

	rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
		rr.DryRun = acpOpts.DryRun
		rr.Confirm = confirm
		rr.OutToStd = true
	})

	if len(args) > 0 {
		rr.GitCmd("status", args...).GitCmd("add", args...)
	} else { // add all changed files
		rr.GitCmd("status").GitCmd("add", ".")
	}

	rr.GitCmd("commit", "-m", acpOpts.message)

	if runPush {
		rr.GitCmd("push")
	}

	return rr.Run()
}
