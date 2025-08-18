package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/errorx"
)

var cmOpts = struct {
	format   cflag.EnumString
	output   string
	withFlag bool
}{
	// output: "stdout",
	format: cflag.NewEnumString("row", "json", "md", "markdown"),
}

// CommandMapCmd command
var CommandMapCmd = &gcli.Command{
	Name:    "cmd-map",
	Aliases: []string{"cmdmap"},
	Desc:    "display all console commands info for kite",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&cmOpts.withFlag, "flags, flag", "output with command options and arguments")
		c.VarOpt(&cmOpts.format, "format", "f", "can be custom the output format")
		c.StrOpt2(&cmOpts.output, "output, o", "custom output, default is stdout", gflag.WithDefault("stdout"))
	},
	Func: func(c *gcli.Command, args []string) error {
		// to json
		if cmOpts.format.String() == "json" {
			// return nil
		}

		return errorx.New("todo")
	},
}

//
// func formatCommandsJSON() {
// 	for name, cmd := range app.Cli.Commands() {
//
// 		for s, command := range cmd.Commands() {
//
// 		}
// 	}
// }
//
// func formatCommandsText() {
// 	for name, cmd := range app.Cli.Commands() {
//
// 	}
// }
//
// func formatCommandsText(c *gcli.Command) {
// 	for s, command := range c.Commands() {
//
// 	}
// }
