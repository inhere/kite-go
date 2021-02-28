package gotool

import "github.com/gookit/gcli/v3"

// AwesomeGo command
var AwesomeGo = &gcli.Command{
	Name: "awesome",
	Desc: "see or search awesome go content",
	Aliases: []string{"awe"},
	Config: func(c *gcli.Command) {
		c.AddArg("keyword", "the keyword for search")
	},
	Func: func(c *gcli.Command, args []string) error {
		return nil
	},
}
