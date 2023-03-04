package textcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/pkg/kiteext"
)

var tsOpts = struct {
	get   gflag.String
	text  string
	match string
}{}

// TextSearchCmd instance
var TextSearchCmd = &gcli.Command{
	Name:    "search",
	Aliases: []string{"find"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {
		c.StrOpt2(&matchOpts.match, "match,m", "set sep char for split input, default is SPACE")
		c.VarOpt2(&matchOpts.get, "get", "get values by indexes, multi by comma")
		c.AddArg("text", "input text contents for process").WithAfterFn(func(a *gflag.CliArg) error {
			matchOpts.text = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := kiteext.ReadContents(matchOpts.text)
		if err != nil {
			return err
		}

		list := strings.Split(src, apputil.ResolveSep(splitOpts.sep))

		fmt.Println(list)
		return nil
	},
}
