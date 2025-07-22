package extcmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// XFileCmd TODO 暂时由 kscript.auto_task_files 替代
var XFileCmd = &gcli.Command{
	Name:    "xfile",
	Aliases: []string{"xrun"},
	Desc: "execute kite xfile command in workdir or parent dir. like makefile, just",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
