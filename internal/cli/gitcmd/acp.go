package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/cmdutil"
	"github.com/inhere/kite/pkg/gitx"
)

var acpOpts = struct {
	gitx.CommonOpts
	template string
	message  string
	notPush  bool
}{}

const acpHelp = `
Commit types:
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

func acpConfigFunc(c *gcli.Command, bindNp bool) {
	acpOpts.BindCommonFlags(c)

	c.StrOpt2(&acpOpts.message, "message,m", "the git commit message", gflag.WithRequired())
	c.StrOpt2(&acpOpts.template, "template,t", "the git commit template")

	if bindNp {
		c.BoolOpt(&acpOpts.notPush, "not-push", "np", false, "dont execute git push")
	}

	c.BindArg(&gcli.CliArg{
		Name:    "files",
		Desc:    "Only add special files. default will add all changed files",
		Arrayed: true,
	})
}

// NewAddCommitPush command
func NewAddCommitPush(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	return &gcli.Command{
		Name: "acp",
		Desc: "run `git add/commit/push` at once command",
		Help: acpHelp,
		Config: func(c *gcli.Command) {
			acpConfigFunc(c, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			// acpOpts.notPush = false
			return acpHandleFunc(c, args, cfgGetter())
		},
	}
}

// NewAddCommitCmd instance
func NewAddCommitCmd(cfgGetter gitx.ConfigProviderFn) *gcli.Command {
	return &gcli.Command{
		Name: "ac",
		Desc: "run `git add/commit` at once command",
		Help: acpHelp,
		Config: func(c *gcli.Command) {
			acpConfigFunc(c, false)
		},
		Func: func(c *gcli.Command, args []string) error {
			return acpHandleFunc(c, args, cfgGetter())
		},
	}
}

func acpHandleFunc(c *gcli.Command, args []string, cfg *gitx.Config) error {
	runPush := !acpOpts.notPush
	confKey := strutil.Join("_", "cmd", cfg.HostType, c.Name)

	dump.NoLoc(app.Cfg().Get(confKey))

	rr := cmdutil.NewRunner(func(rr *cmdutil.Runner) {
		rr.DryRun = acpOpts.DryRun
		rr.Confirm = acpOpts.Confirm
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
