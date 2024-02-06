package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/inhere/kite-go/internal/apputil"
)

// NewTextSearchCmd instance
func NewTextSearchCmd() *gcli.Command {
	var searchOpts = struct {
		get   gflag.String
		text  string
		match string
	}{}

	return &gcli.Command{
		Name:    "search",
		Aliases: []string{"find"},
		Desc:    "search text by pattern, or directly match specified string: date, ip, email, url, phone, etc",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&searchOpts.match, "match,m", "set sep char for split input, default is SPACE")
			c.VarOpt2(&searchOpts.get, "get", "get values by indexes, multi by comma")
			c.AddArg("text", "input text contents for search").WithAfterFn(func(a *gflag.CliArg) error {
				searchOpts.text = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, _ []string) error {
			src, err := apputil.ReadSource(searchOpts.text)
			if err != nil {
				return err
			}

			fmt.Println(src)
			return nil
		},
	}
}

// NewStrMatchCmd create
func NewStrMatchCmd() *gcli.Command {
	var matchOpts = struct {
		get   gflag.String
		text  string
		match string
	}{}

	return &gcli.Command{
		Name:    "match",
		Aliases: []string{"get"},
		Desc:    "simple match and get special part in text",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&matchOpts.match, "match,m", "match and get special part in text. eg: ip, ipv4, ts, date")
			c.VarOpt2(&matchOpts.get, "get", "get values by indexes, multi by comma")
			c.AddArg("text", "input text contents for process").WithAfterFn(func(a *gflag.CliArg) error {
				matchOpts.text = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, _ []string) error {
			src, err := apputil.ReadSource(matchOpts.text)
			if err != nil {
				return err
			}

			// TODO ipv4

			fmt.Println(src)
			return nil
		},
	}
}
