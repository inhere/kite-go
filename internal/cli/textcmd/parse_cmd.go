package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/inhere/kite-go/internal/apputil"
)

// NewTextParseCmd create
func NewTextParseCmd() *gcli.Command {
	parseOpts := struct {
		Text      string `flag:"desc=input text contents for parse;shorts=t"`
		Expr      string `flag:"desc=parse text item by expression pattern;shorts=e"`
		Fields    string `flag:"desc=Set field names for data column, split by ','"`
		GetCol    string `flag:"desc=get column values by indexes, multi by comma, start is 0. eg: 1,5"`
		RowSep    string `flag:"desc=Set row separator;default=NL"`
		ColSep    string `flag:"desc=Set column separator;default=SPACE"`
		ColParser string `flag:"desc=Set column value parser;default=split"`
	}{}

	return &gcli.Command{
		Name:    "parse",
		Desc:    "parse text contents to collect metadata info",
		Aliases: []string{"meta"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&parseOpts)
			c.AddArg("text", "input text contents for process").WithAfterFn(func(a *gflag.CliArg) error {
				parseOpts.Text = a.String()
				return nil
			})
		},
		Examples: ``,
		Help: `
Special keywords:
	- NL: '\n'
	- TAB: '\t'
	- SPACE: ' '
    - COMMA: ','
    - SLASH: '/'
    - BLANK: any blank chars. eg: ' \t\n\r'
`,
		Func: func(c *gcli.Command, _ []string) error {
			src, err := apputil.ReadSource(parseOpts.Text)
			if err != nil {
				return err
			}

			// TODO expr: {}

			fmt.Println(src)
			return nil
		},
	}
}
