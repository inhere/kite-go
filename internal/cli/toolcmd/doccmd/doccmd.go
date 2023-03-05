package doccmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// DocumentCmd instance
var DocumentCmd = &gcli.Command{
	Name:    "doc",
	Desc:    "how to use for common tools or common commands",
	Aliases: []string{"htu", "docs"},
	Subs: []*gcli.Command{
		SearchCmd,
		LinuxCmd,
		InstallCmd,
		TopicListCmd,
		// TODO start web server render markdown pages
	},
	Config: func(c *gcli.Command) {
		c.AddArg("topic", "search document on the topic")
	},
	Func: func(c *gcli.Command, args []string) error {
		topic := c.Arg("topic").String()
		if topic == "" {
			return c.ShowHelp()
		}

		// kite doc linux ls
		// search doc on topic
		return errorx.New("TODO")
	},
}

// SearchCmd instance
var SearchCmd = &gcli.Command{
	Name:    "search",
	Aliases: []string{"find"},
	Desc:    "search document for use linux commands",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
