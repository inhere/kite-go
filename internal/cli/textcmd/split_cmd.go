package textcmd

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/apputil"
)

var splitOpts = struct {
	get  gflag.IntsString
	sep  string
	join string
	text string

	// rowNum int

	count  bool
	noTrim bool
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
			return goutil.OrError(val >= 0, errorx.Raw("get index cannot be < 0"))
		}

		c.StrOpt2(&splitOpts.sep, "sep,s", "set sep char for split input")
		c.StrOpt2(&splitOpts.join, "join", "set join char for build output, default is NL", gflag.WithDefault("NL"))
		c.VarOpt2(&splitOpts.get, "get, i", "get values by indexes, multi by comma")

		c.BoolOpt2(&splitOpts.noTrim, "no-trim", "do not trim input text contents")
		c.BoolOpt2(&splitOpts.count, "count, c", "count item number of split strings")
		c.BoolOpt2(&splitOpts.first, "first, f", "get first part from split strings")
		c.BoolOpt2(&splitOpts.last, "last, l", "get last part from split strings")

		c.AddArg("text", "input text contents for handle").WithAfterFn(func(a *gflag.CliArg) error {
			splitOpts.text = a.String()
			return nil
		})
	},
	Func: strSplitHandle,
	Help: `
### Sep chars:
 NL             - new line
 TAB            - tab char
 SPACE          - space char
 AS,anySpace    - any space chars, like: space, tab, new line
`,
}

var trySeps = [256]uint8{',': 1, ';': 1, ':': 1, '.': 1, '\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func strSplitHandle(_ *gcli.Command, _ []string) error {
	src, err := apputil.ReadSource(splitOpts.text)
	if err != nil {
		return err
	}

	var list = trySplit(src)
	listLen := len(list)
	if listLen == 0 {
		return nil
	}

	var val string
	if splitOpts.first {
		val = list[0]
	} else if splitOpts.last {
		val = list[listLen-1]
	} else {
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

		val = strings.Join(list, joinSep)
	}

	if splitOpts.noTrim {
		fmt.Println(val)
	} else {
		fmt.Print(val)
	}
	return nil
}

func trySplit(s string) []string {
	var list []string
	if splitOpts.sep != "" {
		if splitOpts.sep == "AS" || splitOpts.sep == "anySpace" {
			list = strings.Fields(s)
		} else {
			list = strings.Split(s, apputil.ResolveSep(splitOpts.sep))
		}

		return list
	}

	var sepChar rune
	for i := 0; i < len(s); i++ {
		r := s[i]
		if int(trySeps[r]) == 1 {
			sepChar = rune(r)
			continue
		}
	}

	if sepChar != 0 {
		list = strings.FieldsFunc(s, func(r rune) bool {
			return sepChar == r
		})
	} else {
		list = strings.FieldsFunc(s, func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		})
	}

	return list
}
