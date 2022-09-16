package gitx

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var TagCmd = &gcli.Command{
	Name: "tag",
	Desc: "git tag commands",
	Subs: []*gcli.Command{
		TagCreate,
		TagDelete,
	},
}

var TagCreate = &gcli.Command{
	Name:    "create",
	Aliases: []string{"new"},
	Desc:    "create new tag by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}

var TagDelete = &gcli.Command{
	Name:    "delete",
	Aliases: []string{"del", "rm", "remove"},
	Desc:    "delete exists tags by `git tag`",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
