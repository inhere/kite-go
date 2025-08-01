package fscmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/toolcmd/common"
	"github.com/inhere/kite-go/internal/cli/toolcmd/convcmd"
)

// FsCmd command
var FsCmd = &gcli.Command{
	Name:    "fs",
	Aliases: []string{"file"},
	Desc:    "provide some useful file system commands",
	Subs: []*gcli.Command{
		FileFindCmd,
		ListFilesCmd,
		DeleteCmd,
		RenameCmd,
		NewDirTreeCmd(),
		NewFileCatCmd(),
		NewReplaceCmd(),
		NewTemplateCmd(),
		common.NewQuickOpenCmd(),
		convcmd.NewConvPathSepCmd(),
		// filewatcher.FileWatcher(nil)
		// TODO tree command
	},
}
