package fscmd

import (
	"github.com/gookit/gcli/v3"
)

// FsCmd command
var FsCmd = &gcli.Command{
	Name: "fs",
	// Aliases: []string{"fss"},
	Desc: "provide some useful file system commands",
	Subs: []*gcli.Command{
		FileCatCmd,
		FileFindCmd,
		ListFilesCmd,
		RenameCmd,
		// filewatcher.FileWatcher(nil)
	},
}
