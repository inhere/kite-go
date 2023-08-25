package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/inhere/kite-go/internal/apputil"
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
	Desc:    "search text by pattern, or directly match specified string: date, ip, email, url, phone, etc",
	Config: func(c *gcli.Command) {
		c.StrOpt2(&tsOpts.match, "match,m", "set sep char for split input, default is SPACE")
		c.VarOpt2(&tsOpts.get, "get", "get values by indexes, multi by comma")
		c.AddArg("text", "input text contents for search").WithAfterFn(func(a *gflag.CliArg) error {
			tsOpts.text = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(tsOpts.text)
		if err != nil {
			return err
		}

		fmt.Println(src)
		return nil
	},
}
