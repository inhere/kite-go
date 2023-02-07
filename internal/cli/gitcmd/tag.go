package gitcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// TagCmd instance
var TagCmd = &gcli.Command{
	Name: "tag",
	Desc: "git tag commands",
	Subs: []*gcli.Command{
		TagListCmd,
		TagCreateCmd,
		TagDeleteCmd,
	},
}

// TagListCmd instance
var TagListCmd = &gcli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Desc:    "list tags for the git repository",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}

// TagCreateCmd instance
var TagCreateCmd = &gcli.Command{
	Name:    "create",
	Aliases: []string{"new"},
	Desc:    "create new tag by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}

// TagDeleteCmd instance
var TagDeleteCmd = &gcli.Command{
	Name:    "delete",
	Aliases: []string{"del", "rm", "remove"},
	Desc:    "delete exists tags by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
