package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// MathCalcCmd command
var MathCalcCmd = &gcli.Command{
	Name:    "math",
	Aliases: []string{"calc"},
	Desc:    "list the jump storage data in local",
	Config: func(c *gcli.Command) {
		// random string(number,alpha,), int(range)
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
