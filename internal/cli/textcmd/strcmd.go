package textcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite/internal/apputil"
	"github.com/inhere/kite/pkg/kiteext"
)

// TextOperateCmd instance
var TextOperateCmd = &gcli.Command{
	Name:    "text",
	Desc:    "useful commands for handle string text",
	Aliases: []string{"str", "string"},
	Subs: []*gcli.Command{
		StrCountCmd,
		StrSplitCmd,
		StrMatchCmd,
		TextSearchCmd,
		// TODO
	},
}

// StrCountCmd instance
var StrCountCmd = &gcli.Command{
	Name:    "length",
	Aliases: []string{"len", "count"},
	Desc:    "count input string length",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

var splitOpts = struct {
	get   gflag.IntsString
	sep   string
	join  string
	text  string
	count bool
	// quick fetch
	first, last bool
}{}

// StrSplitCmd instance
var StrSplitCmd = &gcli.Command{
	Name: "split",
	// Aliases: []string{"len", "count"},
	Desc: "split input text to multi parts, then fetch or joins",
	Config: func(c *gcli.Command) {
		splitOpts.get.ValueFn = func(val int) error {
			return goutil.OrError(val < 0, errorx.Rawf("get index cannot be < 0"))
		}

		c.StrOpt2(&splitOpts.sep, "sep,s", "set sep char for split input, default is SPACE", gflag.WithDefault("SPACE"))
		c.StrOpt2(&splitOpts.join, "join", "set join char for build output, default is NL", gflag.WithDefault("NL"))
		c.VarOpt2(&splitOpts.get, "get", "get values by indexes, multi by comma")

		c.BoolOpt2(&splitOpts.count, "count, c", "get first part from split strings")
		c.BoolOpt2(&splitOpts.first, "first, f", "get first part from split strings")
		c.BoolOpt2(&splitOpts.last, "last, l", "get last part from split strings")

		c.AddArg("text", "input text contents for handle").WithAfterFn(func(a *gflag.CliArg) error {
			splitOpts.text = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := kiteext.ReadContents(splitOpts.text)
		if err != nil {
			return err
		}

		list := strings.Split(src, apputil.ResolveSep(splitOpts.sep))
		listLen := len(list)
		if listLen == 0 {
			return nil
		}

		if splitOpts.first {
			fmt.Println(list[0])
			return nil
		}
		if splitOpts.last {
			fmt.Println(list[listLen-1])
			return nil
		}

		joinSep := apputil.ResolveSep(splitOpts.join)

		if indexes := splitOpts.get.Ints(); len(indexes) > 0 {
			newList := make([]string, 0, len(indexes))
			for _, i := range indexes {
				if i < listLen {
					newList = append(newList, list[i])
				}
			}
			list = newList
		}

		fmt.Println(strings.Join(list, joinSep))
		return nil
	},
}