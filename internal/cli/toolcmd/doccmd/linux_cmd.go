package doccmd

import (
	"errors"

	"github.com/gookit/gcli/v3"
)

// LinuxCmd instance
// https://github.com/jaywcjlove/linux-command/tree/master/command  linux commands zh-CN documents
//   - raw contents eg: https://raw.githubusercontent.com/jaywcjlove/linux-command/master/command/accept.md
//     结构使用 https://raw.githubusercontent.com/jaywcjlove/linux-command/master/dist/data.json
var LinuxCmd = &gcli.Command{
	Name:    "linux",
	Aliases: []string{"lin", "linux-cmd"},
	Desc:    "document for use linux commands",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}
