package mkdown

import "github.com/gookit/gcli/v3"

var MkDownCmd = &gcli.Command{
	Name: "mkdown",
	Desc: "some tool for markdown",
	Aliases: []string{"md", "markdown"},
	Subs: []*gcli.Command{
		Markdown2HTML, Markdown2SQL,
	},
}
