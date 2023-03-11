package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/byteutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/internal/apputil"
)

// ReadmeCmd instance
var ReadmeCmd = &gcli.Command{
	Name:    "readme",
	Aliases: []string{"doc", "docs"},
	Desc:    "show readme docs for kite app",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		str := byteutil.SafeString(kite.EmbedFs.ReadFile("README.md"))
		// fmt.Println(byteutil.SafeString(kite.EmbedFs.ReadFile("README.md")))

		return apputil.RenderContents(str, "markdown", "github")
	},
}
