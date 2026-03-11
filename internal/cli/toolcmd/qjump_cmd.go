package toolcmd

import (
	"strings"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/x/stdio"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/quickjump"
)

// QuickJumpCmd command
var QuickJumpCmd = &gcli.Command{
	Name:    "jump",
	Aliases: []string{"goto"},
	Desc:    "Jump helps you navigate faster by your history.",
	Subs: []*gcli.Command{
		AutoJumpListCmd,
		AutoJumpShellCmd,
		AutoJumpMatchCmd,
		AutoJumpGetCmd,
		AutoJumpSetCmd,
		AutoJumpChdirCmd,
		QuickJumpCleanCmd,
	},
	Config: func(c *gcli.Command) {

	},
}

// AutoJumpListCmd command
var AutoJumpListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list the jump storage data in local",
	Config: func(c *gcli.Command) {
		c.AddArg("type", "the jump info type name. allow: prev,last,history,all")
	},
	Func: func(c *gcli.Command, _ []string) error {
		colorp.Infof("Jump storage datafile: %s\n", app.QJump.Datafile())
		colorp.Greenln("Jump storage data in local:")

		dp := dump.NewWithOptions(dump.WithoutPosition(), dump.SkipPrivate())
		dp.Println(app.QJump.Metadata)
		show.MList(app.QJump.Metadata)
		return nil
	},
}

var jsOpts = struct {
	Bind string `flag:"set the bind func name;false;jump;func"`
}{}

// AutoJumpShellCmd command
var AutoJumpShellCmd = &gcli.Command{
	Name:    "shell",
	Aliases: []string{"active", "script"},
	Desc:    "Generate shell script for give shell env name.",
	Help: `
  Enable quick jump for bash(add to <mga>~/.bashrc</>):
    # shell func is: jump
    <mga>eval "$(kite tool jump shell bash)"</>

Enable quick jump for zsh(add to <mga>~/.zshrc</>):
    # shell func is: jump
    <mga>eval "$(kite tool jump shell zsh)"</>
    # set the bind func name is: j
    <mga>eval "$(kite tool jump shell --bind j zsh)"</>

Enable quick jump for pwsh(add to <mga>$PROFILE</>):
    # jump func is: j
    <mga>kite tool jump shell --bind j pwsh | Out-String | Invoke-Expression</>

`,
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&jsOpts, gflag.TagRuleSimple)
		c.AddArg("shell", "The shell name. allows: bash, zsh, fish, pwsh.")
	},
	Func: func(c *gcli.Command, _ []string) error {
		shellName := c.Arg("shell").String()
		if shellName == "" {
			shellName = sysutil.CurrentShell(true)
		}

		if !quickjump.IsSupported(shellName) {
			colorp.Redf("The shell %q is not supported yet!\n", shellName)
			return nil
		}

		script, err := quickjump.GenScript(shellName, jsOpts.Bind)
		if err != nil {
			return err
		}

		stdio.WriteString(script)
		return nil
	},
}

const (
	ScopeAll = iota
	ScopeNamed
	ScopeHistory
)

var ajmOpts = struct {
	OnlyPath bool `flag:"only-path;only show the path, dont with name"`
	Limit    int  `flag:"limit the match count;false;10"`
	Scope    int  `flag:"scopes for search, allow: 0=match all,1=from named,2=from history;false;0"`
}{}

// AutoJumpMatchCmd command
var AutoJumpMatchCmd = &gcli.Command{
	Name:    "search",
	Aliases: []string{"hint", "match"},
	Desc:    "Match directory paths by given keywords",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&ajmOpts, gflag.TagRuleSimple)
		c.AddArg("keywords", "The keywords to match, allow limit by multi", false, true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		var results []string
		keywords := c.Arg("keywords").Strings()
		keywords = app.QJump.FormatKeywords(keywords)

		// from named
		if ajmOpts.Scope == ScopeNamed {
			results = app.QJump.SearchNamed(keywords, ajmOpts.Limit, !ajmOpts.OnlyPath)
		} else if ajmOpts.Scope == ScopeHistory {
			results = app.QJump.SearchHistory(keywords, ajmOpts.Limit)
		} else {
			results = app.QJump.Search(keywords, ajmOpts.Limit, !ajmOpts.OnlyPath)
		}

		var sb strings.Builder
		matchNum := len(results)
		app.L.Infof("input search keywords %v, search results: %v", keywords, results)

		for i, dirPath := range results {
			sb.WriteString(dirPath)
			if matchNum != i+1 {
				sb.WriteByte('\n')
			}
		}

		stdio.WriteString(sb.String())
		return nil
	},
}

var ajgOpts = struct {
	Quiet bool `flag:"quiet mode, dont report error on missing"`
}{}

// AutoJumpGetCmd command
var AutoJumpGetCmd = &gcli.Command{
	Name:    "get",
	Aliases: []string{"path"},
	Desc:    "Get the real directory path by given name.",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&ajgOpts, gflag.TagRuleSimple)
		c.AddArg("keywords", "The target directory name or path or keywords.\n Start with ^ for exclude", true, true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		keywords := c.Arg("keywords").Strings()
		dirPath := app.QJump.CheckOrMatch(keywords)

		if !ajgOpts.Quiet && len(dirPath) == 0 {
			return c.NewErrf("not found path by keywords: %v", keywords)
		}

		app.L.Infof("input keywords: %v, match dirPath: %s", keywords, dirPath)
		stdio.WriteString(dirPath)
		return nil
	},
}

// AutoJumpSetCmd command
var AutoJumpSetCmd = &gcli.Command{
	Name:    "set",
	Aliases: []string{"add"},
	Desc:    "Set the name to real directory path mapping",
	Config: func(c *gcli.Command) {
		c.AddArg("name", "The name of the directory path", true)
		c.AddArg("path", "The real directory path", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		name := c.Arg("name").String()
		path := c.Arg("path").String()

		if app.QJump.AddNamed(name, path) {
			colorp.Successf("Set jump name %q to path %q\n", name, path)
		} else {
			colorp.Warnln("Set jump name %q to path %q failed", name, path)
		}

		return nil
	},
}

var ajcOpts = struct {
	Quiet bool `flag:"Quiet to add the path to history"`
}{}

// AutoJumpChdirCmd command
var AutoJumpChdirCmd = &gcli.Command{
	Name:    "chdir",
	Aliases: []string{"into", "to"},
	Desc:    "add directory path to history, by the jump dir hooks.",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&ajcOpts, gflag.TagRuleSimple)
		c.AddArg("path", "The real directory path. if empty, use current workdir")
	},
	Func: func(c *gcli.Command, _ []string) error {
		path := c.Arg("path").String()
		if len(path) == 0 {
			return c.NewErrf("path is empty or invalid")
		}

		realPath, ok := app.QJump.AddHistory(path)
		if ok {
			app.L.Infof("jump to path %q success", realPath)
			if !ajcOpts.Quiet {
				colorp.Successf("Into %q\n", realPath)
			}
		} else {
			colorp.Warnf("Invalid path %q\n", realPath)
		}

		return nil
	},
}

// QuickJumpCleanCmd command
var QuickJumpCleanCmd = &gcli.Command{
	Name:    "clean",
	Aliases: []string{"clear"},
	Desc:    "clean invalid directory paths from history",
	Config: func(c *gcli.Command) {
		c.AddArg("path", "The history directory path. if empty, clean all invalid dirs")
	},
	Func: func(c *gcli.Command, args []string) error {
		ss := app.QJump.CleanHistories()

		if len(ss) > 0 {
			show.AList("Cleaned invalid paths", ss)
		} else {
			colorp.Infoln("No invalid paths to clean")
		}
		return nil
	},
}
