package toolcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/app"
)

// AutoJumpCmd command
var AutoJumpCmd = &gcli.Command{
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
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errorx.New("TODO")
	// },
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
		show.MList(app.QJump)
		return nil
	},
}

var jsOpts = struct {
	Bind string `flag:"set the bind func name;false;$SHELL"`
}{}

// AutoJumpShellCmd command
var AutoJumpShellCmd = &gcli.Command{
	Name:    "shell",
	Aliases: []string{"active"},
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
`,
	Config: func(c *gcli.Command) {
		c.UseSimpleRule()
		c.MustFromStruct(&jsOpts)
		c.AddArg("shell", "The shell name. eg: bash, zsh, fish, etc.")
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// AutoJumpMatchCmd command
var AutoJumpMatchCmd = &gcli.Command{
	Name:    "match",
	Aliases: []string{"hint", "search"},
	Desc:    "Match directory paths by given keywords",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// AutoJumpGetCmd command
var AutoJumpGetCmd = &gcli.Command{
	Name:    "get",
	Aliases: []string{"path"},
	Desc:    "Get the real directory path by given name.",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
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
			colorp.Successf("Set jump name %q to path %q success\n", name, path)
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
	Desc:    "add target directory path to history, by the jump dir hooks.",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&ajcOpts, gflag.TagRuleSimple)
		c.AddArg("path", "The real directory path", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		path := c.Arg("path").String()

		if app.QJump.AddHistory(path) {

		}

		return errorx.New("TODO")
	},
}
