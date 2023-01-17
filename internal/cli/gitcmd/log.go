package gitcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
)

var (
	logOpts = struct {
		Abbrev    bool `flag:"Only display the abbrev commit ID"`
		NoColor   bool `flag:"Dont use color render git output"`
		NoMerges  bool `flag:"No contains merge request logs"`
		MaxCommit int  `flag:"Max display how many commits;;15"`
		Format    string
		RepoDir   string      `flag:"repo directory for run git log, default is work dir"`
		Logfile   string      `flag:"export changelog message to file"`
		Exclude   gcli.String `flag:"exclude contains given sub-string. multi by comma split."`
	}{}

	ShowLog = &gcli.Command{
		Name: "log",
		Desc: "display recently git commits information by `git log`",
		// Aliases: []string{"cl", "clog", "changelog"},
		Config: func(c *gcli.Command) {
			c.UseSimpleRule()

			// goutil.PanicIfErr(c.FromStruct(&logOpts))
			goutil.PanicIfErr(c.FromStruct(&logOpts))

			c.StrOpt(&logOpts.Format, "format", "", "",
				"The git log option '--pretty' value.\n"+
					"can be one of oneline, short, medium, full, fuller, reference, email, raw, format:string and tformat:string.",
			)

			c.AddArg("maxCommit", "Max display how many commits")
		},
		Func: func(c *gcli.Command, args []string) error {
			maxNum := c.Arg("maxCommit").Int()
			// git log --color --graph --pretty=format:'%Cred%h%Creset:%C(ul yellow)%d%Creset %s (%Cgreen%cr%Creset, %C(bold blue)%an%Creset)' --abbrev-commit -10
			gitLog := gitw.New("log", "--graph")
			gitLog.OnBeforeExec(gitw.PrintCmdline)

			gitLog.Argf("-%d", maxNum)
			gitLog.ArgIf("--color", !logOpts.NoColor)
			gitLog.ArgIf("--no-merges", logOpts.NoMerges)
			gitLog.ArgIf("--abbrev-commit", logOpts.Abbrev)
			gitLog.Add(`--pretty=format:%Cred%h%Creset:%C(ul yellow)%d%Creset %s (%Cgreen%cr%Creset, %C(bold blue)%an%Creset)`)

			// dump.P(logOpts, maxNum)

			return gitLog.Run()
		},
	}
)
