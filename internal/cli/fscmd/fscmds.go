package fscmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// FsCmd command
var FsCmd = &gcli.Command{
	Name: "fs",
	// Aliases: []string{"fss"},
	Desc: "provide some useful file system commands",
	Subs: []*gcli.Command{
		FileCat,
		FileFinder,
		// filewatcher.FileWatcher(nil)
	},
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
