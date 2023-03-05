package mdcmd

import "github.com/gookit/gcli/v3"

// https://github.com/MichaelMure/go-term-markdown

// MkDownCmd instance
var MkDownCmd = &gcli.Command{
	Name:    "mkdown",
	Desc:    "some tool for markdown",
	Aliases: []string{"md", "mdown", "markdown"},
	Subs: []*gcli.Command{
		MdRenderCmd,
		Markdown2HTML,
		Markdown2SQL,
		MarkdownServeCmd,
	},
}
