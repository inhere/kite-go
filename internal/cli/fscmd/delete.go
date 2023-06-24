package fscmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
)

// DeleteCmd instance
var DeleteCmd = &gcli.Command{
	Name:    "delete",
	Desc:    "delete files by glob pattern",
	Aliases: []string{"del", "rm"},
	Config: func(c *gcli.Command) {
		// TODO regex: from (\w+)_(\w+) to $1_new_$2
	},
	Func: func(c *gcli.Command, _ []string) error {
		colorp.Infoln("TIP: please use the: kite fs find command with option '--del' to delete find files")
		return nil
	},
}
