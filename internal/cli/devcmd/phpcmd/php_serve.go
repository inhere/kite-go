package phpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// PhpServeConf struct
// run: php -S 127.0.0.1:8080 -t web web/index.php
type PhpServeConf struct {
	Root   string
	Entry  string
	PhpBin string
	Addr   string
}

// PhpServeCmd instance
var PhpServeCmd = &gcli.Command{
	Name: "serve",
	Desc: "start an php development server",
	Func: func(c *gcli.Command, args []string) error {
		cmd := cmdr.NewCmd("php", "-v")

		return cmd.Run()
	},
}
