package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
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
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// AutoJumpListCmd command
var AutoJumpListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list the jump storage data in local",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// AutoJumpShellCmd command
var AutoJumpShellCmd = &gcli.Command{
	Name:    "shell",
	Aliases: []string{"active"},
	Desc:    "Generate shell script for give shell env name.",
	Help: `
  quick jump for bash(add to ~/.bashrc):
      # shell func is: jump
      eval "$(kite jump shell bash)"

  quick jump for zsh(add to ~/.zshrc):
      # shell func is: jump
      eval "$(kite jump shell zsh)"
      # set the bind func name is: j
      eval "$(kite jump shell zsh --bind j)"
`,
	Config: func(c *gcli.Command) {

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
	Aliases: []string{"cd"},
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

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// AutoJumpChdirCmd command
var AutoJumpChdirCmd = &gcli.Command{
	Name:    "chdir",
	Aliases: []string{"into"},
	Desc:    "record target directory path, by the jump dir hooks.",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
