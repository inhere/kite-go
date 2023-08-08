package textcmd

import (
	"regexp"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/pkg/kautorw"
)

var trOpts = struct {
	From string `flag:"replace text from"`
	To   string `flag:"replace text to"`
	// Expr like /FROM/TO/
	Expr  string `flag:"quickly replace text by rule expression. FORMAT: /FROM/TO/"`
	Write bool   `flag:"write result to src file, on input is filepath;;;w"`
	Regex bool   `flag:"replace text by regex expression, mark --from and --to as regex pattern;;;r"`
	// text string
	text string
}{}

// TextReplaceCmd instance
var TextReplaceCmd = &gcli.Command{
	Name:    "replace",
	Aliases: []string{"repl", "rpl"},
	Desc:    "simple and quickly replace text contents",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&trOpts, gflag.TagRuleSimple)
		c.AddArg("text", "input text contents for process. allow @c,@FILE").WithAfterFn(func(a *gflag.CliArg) error {
			trOpts.text = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(tsOpts.text)
		if err != nil {
			return err
		}

		if trOpts.Expr != "" {
			trOpts.From, trOpts.To = strutil.QuietCut(strings.Trim(trOpts.Expr, "/"), "/")
		}

		var ret string
		if trOpts.Regex {
			reg := regexp.MustCompile(trOpts.From)
			ret = reg.ReplaceAllString(src, trOpts.To)
		} else {
			ret = strings.ReplaceAll(src, trOpts.From, trOpts.To)
		}

		sw := kautorw.NewSourceWriter("")
		sw.SetSrcFile(trOpts.text)

		if trOpts.Write {
			sw.WithDst("@src")
			if !sw.HasSrcFile() {
				return c.NewErrf("with option --write, but input is not a file")
			}
		}

		return sw.WriteString(ret)
	},
}
