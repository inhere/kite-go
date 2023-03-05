package doccmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// Topic struct
type Topic struct {
	Name string
	Desc string
	Type string // file, git, dir
	Path string
}

// InstallCmd instance
var InstallCmd = &gcli.Command{
	Name:    "install",
	Aliases: []string{"ins", "add"},
	Desc:    "install new documents topic from git repository(eg: github)",
	Config: func(c *gcli.Command) {
		// type: git
		//
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// TopicListCmd instance
var TopicListCmd = &gcli.Command{
	Name:    "list-topics",
	Aliases: []string{"list-topic", "lt"},
	Desc:    "list installed document topics",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
