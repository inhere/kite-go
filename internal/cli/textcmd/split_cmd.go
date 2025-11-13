package textcmd

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/inhere/kite-go/internal/apputil"
)

var splitOpts = struct {
	get  gflag.IntsString
	sep  string // 分隔符
	join string // 构建输出结果时使用的分隔符
	text string

	// rowNum int

	count  bool
	noTrim bool
	inline bool
	// 不删除 item 前后的空白
	noTrimItems bool
	// quick fetch
	first, last bool
	// 使用自定义模板格式化输出 eg: "new $1 format $2"
	// - 参数是 item index: $0, $1, ... $N
	format string
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

		c.StrOpt2(&splitOpts.sep, "sep,s", "set sep char for split input. if not set, will auto detect available char.")
		c.VarOpt2(&splitOpts.get, "get, i", "get values by indexes, multi by comma")

		c.BoolOpt2(&splitOpts.inline, "inline", "inline output result, dont end with NL")
		c.StrOpt2(&splitOpts.join, "join", "set join char for build output", gflag.WithDefault("NL"))
		c.StrOpt2(&splitOpts.format, "format", "format output result, use $0, $1, ... $N")

		c.BoolOpt2(&splitOpts.noTrim, "no-trim", "do not trim input text contents")
		c.BoolOpt2(&splitOpts.noTrimItems, "no-trim-items", "do not trim each item contents")

		c.BoolOpt2(&splitOpts.count, "count, c", "count item number of split strings")
		c.BoolOpt2(&splitOpts.first, "first, 0", "get first part from split strings")
		c.BoolOpt2(&splitOpts.last, "last, l", "get last part from split strings")

		c.AddArg("text", "input text contents for handle. allow use: @c, @i").WithAfterFn(func(a *gflag.CliArg) error {
			splitOpts.text = a.String()
			return nil
		})
	},
	Func: strSplitHandle,
	Help: `
<green>### Sep Chars</>:
 NL             - new line
 TAB            - tab char
 SPACE          - space char
 AS,anySpace    - any space chars, like: space, tab, new line
`,
	Examples: `
# custom output format
{$binWithCmd} --last 'a b c' # Output: c
{$binWithCmd} --first 'a b c' # Output: a
{$binWithCmd} --format '$0,$2' 'a b c' # Output: a,c
`,
}

var trySeps = [256]uint8{',': 1, ';': 1, ':': 1, '.': 1, '\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
var fmtRpl = textutil.NewVarReplacer("$")

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

	// 只显示统计
	if splitOpts.count {
		fmt.Println(listLen)
		return nil
	}

	// 去除每个item的空白
	if !splitOpts.noTrimItems {
		for i := 0; i < listLen; i++ {
			list[i] = strings.TrimSpace(list[i])
		}
	}

	var val string
	if splitOpts.first {
		val = list[0]
	} else if splitOpts.last {
		val = list[listLen-1]
	} else if splitOpts.format != "" {
		val = fmtRpl.RenderSimple(splitOpts.format, listToVarMap(list))
	} else {
		// 获取指定索引的值
		if indexes := splitOpts.get.Ints(); len(indexes) > 0 {
			newList := make([]string, 0, len(indexes))
			for _, i := range indexes {
				if i < listLen {
					newList = append(newList, list[i])
				}
			}
			list = newList
		}

		joinSep := apputil.ResolveSep(splitOpts.join)
		val = strings.Join(list, joinSep)
	}

	if splitOpts.inline {
		fmt.Print(val)
	} else {
		fmt.Println(val)
	}
	return nil
}

func listToVarMap(list []string) map[string]string {
	mp := make(map[string]string, len(list))
	for i, v := range list {
		mp[strutil.SafeString(i)] = v
	}
	return mp
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

	// 自动检测分隔符
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
