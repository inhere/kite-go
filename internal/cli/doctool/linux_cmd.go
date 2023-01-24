package doctool

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// LinuxCmd instance
var LinuxCmd = &gcli.Command{
	Name:    "linux",
	Aliases: []string{"lin", "linux-cmd"},
	Desc:    "document for use linux commands",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
