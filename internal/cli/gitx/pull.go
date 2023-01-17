package gitx

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
)

// UpdateCmd command
var UpdateCmd = &gcli.Command{
	Name:    "pull",
	Desc:    "Update codes from git remote repositories",
	Aliases: []string{"pul", "pl"},
	Config: func(c *gcli.Command) {
		bindCommonFlags(c)
	},
	Func: func(c *gcli.Command, args []string) error {
		pull := gitw.NewWithArgs("pull", args...)
		pull.WithWorkDir(workdir)

		return pull.Run()
	},
}
