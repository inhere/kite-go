package fscmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/textcmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd/convcmd"
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
		convcmd.NewConvPathSepCmd(),
		textcmd.NewTemplateCmd(),
		// filewatcher.FileWatcher(nil)
	},
}
