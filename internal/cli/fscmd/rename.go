package fscmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
)

// RenameCmd instance
var RenameCmd = &gcli.Command{
	Name: "rename",
	Desc: "rename files by glob or regexp pattern",
	Config: func(c *gcli.Command) {
		// TODO regex: from (\w+)_(\w+) to $1_new_$2
	},
	Func: func(c *gcli.Command, _ []string) error {
		colorp.Infoln("TIP: please use the: kite fs find command with option '--exec' to rename find files")
		return nil
	},
}
