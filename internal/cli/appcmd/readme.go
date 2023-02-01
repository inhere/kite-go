package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/byteutil"
	"github.com/inhere/kite"
)

// ReadmeCmd instance
var ReadmeCmd = &gcli.Command{
	Name:    "readme",
	Aliases: []string{"doc", "docs"},
	Desc:    "show readme docs for kite app",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		fmt.Println(byteutil.SafeString(kite.EmbedFs.ReadFile("README.md")))
		return nil
	},
}
