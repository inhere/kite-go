package gotool

import "github.com/gookit/gcli/v3"

var awesomeEnUrl = "https://github.com/avelino/awesome-go"
var awesomeCnUrl = "https://github.com/yinggaozhen/awesome-go-cn"

// AwesomeGo command
var AwesomeGo = &gcli.Command{
	Name:    "awesome",
	Desc:    "view or search package on awesome go content",
	Aliases: []string{"awe"},
	Config: func(c *gcli.Command) {
		c.AddArg("keyword", "the keyword for search")
	},
	Func: func(c *gcli.Command, args []string) error {
		return nil
	},
}
