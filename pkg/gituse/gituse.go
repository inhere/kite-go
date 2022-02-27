package gituse

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

var (
	DryRun  bool
	Workdir string
)

func BindCommonFlags(c *gcli.Command) {
	c.BoolOpt(&DryRun, "dry-run", "", false, "run workflow, but dont real execute command")
	c.StrOpt(&Workdir, "workdir", "w", "", "the command workdir path")
}

// OpenRemoteRepo address
var OpenRemoteRepo = &gcli.Command{
	Name: "open",
	Desc: "open the git remote repo address",
	Func: func(c *gcli.Command, args []string) error {
		return errors.New("TODO")
	},
}
