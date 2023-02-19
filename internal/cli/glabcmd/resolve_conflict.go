package glabcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite/pkg/gitx"
)

var rcOpts gitx.CommonOpts

// ResolveConflictCmd instance
var ResolveConflictCmd = &gcli.Command{
	Name: "resolve",
	Desc: "Resolve conflicts preparing for current git branch.",
	Help: `Workflow:
1. will checkout to <cyan>branch</cyan>
2. will update code by <cyan>git pull</cyan>
3. update the <cyan>branch</cyan> codes from source repository
4. merge current-branch codes from source repository
5. please resolve conflicts by tools or manual
`,
	Aliases: []string{"rc"},
	Config: func(c *gcli.Command) {
		rcOpts.BindCommonFlags(c)

		c.AddArg("branch", "The conflicts target branch name. eg: qa, pre, master", true)
	},
	Func: func(c *gcli.Command, args []string) error {

		return errorx.New("TODO")
	},
}
