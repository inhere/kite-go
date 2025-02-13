package common

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/app"
)

// NewQuickOpenCmd command
func NewQuickOpenCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "open",
		Aliases: []string{"open-file"},
		Desc:    "open input file or dir or remote URL address",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "bin name or URL address", true)
		},
		Func: func(c *gcli.Command, _ []string) error {
			name := c.Arg("name").String()

			var dstFile = name
			if strings.Contains(name, "/") {
				// special github url
				if strings.HasPrefix(name, gitw.GitHubHost) {
					dstFile = "https://" + name
					// } else if fsutil.PathExists(name) {
					// 	// nothing ...
					// } else if validate.IsURL(name) {
				}
			} else if app.OpenMap.HasAlias(name) {
				dstFile = app.OpenMap.ResolveAlias(name)
			}

			c.Infoln("Will Open the:", dstFile)
			return sysutil.Open(dstFile)
		},
	}
}
