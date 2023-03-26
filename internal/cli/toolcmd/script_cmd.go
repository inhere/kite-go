package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
)

// ScriptCmd command
var ScriptCmd = &gcli.Command{
	Name:    "script",
	Aliases: []string{"scripts"},
	Desc:    "display or search available scripts or script-file",
	Config: func(c *gcli.Command) {
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
