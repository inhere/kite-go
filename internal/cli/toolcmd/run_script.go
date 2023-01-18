package toolcmd

import "github.com/gookit/gcli/v3"

var runOpts = struct {
	listAll bool
	search  bool
}{}
var RunScripts = &gcli.Command{
	Name:    "run",
	Desc:    "run an custom script command in the `scripts`",
	Aliases: []string{"exec", "script"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(
			&runOpts.listAll,
			"list", "l", false, "List information for all scripts or one script",
		)
		c.BoolOpt(
			&runOpts.search,
			"search", "s", false, "Display all matched scripts by the input name",
		)

		c.AddArg("name", "The script name for execute or display")
	},
	Func: func(c *gcli.Command, args []string) error {
		c.Infoln("TODO")
		return nil
	},
	Help: `
Can use '$@' '$?' at script line. will auto replace to input arguments
examples:

  #yaml
  st: git status
  co: git checkout $@
  br: git branch $?
`,
}
