package fscmd

import (
	"github.com/gookit/gcli/v3"
)

// FsCmd command
var FsCmd = &gcli.Command{
	Name:    "fs",
	Aliases: []string{"file"},
	Desc:    "provide some useful file system commands",
	Subs: []*gcli.Command{
		NewFileCatCmd(),
		FileFindCmd,
		ListFilesCmd,
		RenameCmd,
		DeleteCmd,
		// filewatcher.FileWatcher(nil)
	},
}
