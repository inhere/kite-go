package gitcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw"
	"github.com/gookit/gitw/gitutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/internal/biz/cmdbiz"
	"github.com/inhere/kite/pkg/cmdutil"
	"github.com/inhere/kite/pkg/gitx"
)

type acpOptModel struct {
	cmdbiz.CommonOpts
	template string
	message  string

	notPush    bool
	noTemplate bool
	autoEmoji  bool
	autoType   bool
	autoSign   bool
}

func (m *acpOptModel) buildMsg(tpl, brName string) string {
	if m.noTemplate {
		return m.message
	}

	tpl = strutil.OrElse(m.template, tpl)
	if tpl == "" {
		return m.message
	}

	msgVar := "{message}"
	if !strings.Contains(tpl, msgVar) {
		tpl = fmt.Sprintf("%s %s", tpl, msgVar)
	}

	topics := gitutil.ParseCommitTopic(m.message)

	vars := map[string]any{
		"branch":  brName,
		"message": m.message,
		"topic":   arrutil.Strings(topics).First(),
		"emoji":   "", // TODO
	}

	return textutil.ReplaceVars(tpl, vars, appconst.VarFormat)
}

const acpHelp = `
<b>Commit types</>:
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

<b>Template variables</>:
> variables can use for message template

emoji   - auto add git-emoji
topic   - auto add commit topic type
branch  - current branch name
message - input commit message

Examples:

{$fullCmd} -t '{emoji} {message}' -m "fix: fix an error"
> Will run: <cyan>git -m ":bug: fix: fix an error"</>
`

var acpOpts = acpOptModel{}

func acpConfigFunc(c *gcli.Command, bindNp bool) {
	acpOpts.BindCommonFlags(c)

	c.StrOpt2(&acpOpts.message, "message,m", "the git commit message", gflag.WithRequired())
	c.StrOpt2(&acpOpts.template, "template,t", "the git commit template")
	c.BoolOpt(&acpOpts.noTemplate, "no-template", "nt", false, "disable the commit template")
	c.BoolOpt(&acpOpts.autoEmoji, "auto-emoji", "ae,emoji", false, "auto prepend git-emoji to commit template")

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
	confKey := apputil.CmdConfigKey(cfg.HostType, c.Name)
	cmdConf := app.Cfg().SubDataMap(confKey)

	if !cmdConf.IsEmtpy() {
		show.AList("Command settings", cmdConf)
	}

	gr := gitw.NewRepo(acpOpts.Workdir)
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

	message := acpOpts.buildMsg(cmdConf.Str("template"), gr.CurBranchName())
	rr.GitCmd("commit", "-m", message)

	if runPush {
		rr.GitCmd("push")
	}
	return rr.Run()
}
