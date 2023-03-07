package toolcmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite/pkg/kiteext"
)

// ConvBaseCmd command
// Base
// Binary
// decimal
// Base 8
var ConvBaseCmd = &gcli.Command{
	Name:    "conv-base",
	Aliases: []string{"base", "cb"},
	Desc:    "list the jump storage data in local",
	Config: func(c *gcli.Command) {
		// random string(number,alpha,), int(range)
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

var timeRegex = regexp.MustCompile(`1\d{9}`)

// Time2dateCmd instance
var Time2dateCmd = &gcli.Command{
	Name:    "ts2date",
	Aliases: []string{"ts", "t2d", "t2date"},
	Desc:    "Quick convert all timestamp number to datetime",
	Config: func(c *gcli.Command) {
		c.AddArg("input", "want parsed input contents", true, true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		var txt string
		ss := c.Arg("input").Strings()
		if len(ss) > 1 {
			txt = strings.Join(ss, " ")
		}

		txt, err := kiteext.ReadContents(ss[0])
		if err != nil {
			return err
		}

		c.Infoln("Input Contents:")
		fmt.Println(txt + "\n")

		times := timeRegex.FindAllString(txt, -1)
		if len(times) == 0 {
			return errorx.Raw("not found any timestamps")
		}

		mp := make(map[string]string, len(times))
		for _, timeVal := range times {
			mp[timeVal] = timex.FormatUnix(mathutil.SafeInt64(timeVal))
		}

		show.AList("Matched timestamps", mp, func(opts *show.ListOption) {
			opts.SepChar = "  =>  "
		})
		return nil
	},
}
