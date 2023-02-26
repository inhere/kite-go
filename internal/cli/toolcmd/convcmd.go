package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// ConvBaseCmd command
// Base
// Binary
// decimal
// Base 8
var ConvBaseCmd = &gcli.Command{
	Name:    "conv-base",
	Aliases: []string{"base", "cb"},
	Desc:    "list the jump storage data in local",
	Config: func(c *gcli.Command) {
		// random string(number,alpha,), int(range)
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
