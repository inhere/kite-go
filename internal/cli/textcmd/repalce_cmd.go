package textcmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/pkg/kautorw"
)

var trOpts = struct {
	From string `flag:"desc=replace text from;shorts=f"`
	To   string `flag:"desc=replace text to;shorts=t"`
	// Expr like /FROM/TO/
	Expr  string `flag:"desc=quickly replace text by rule expression. FORMAT: /FROM/TO/"`
	Write bool   `flag:"desc=write result to src file, on input is filepath;shorts=w"`
	Regex bool   `flag:"desc=replace text by regex expression, mark --from and --to as regex pattern;shorts=r"`
	// text string
	text string
}{}

// TextReplaceCmd instance
var TextReplaceCmd = &gcli.Command{
	Name:    "replace",
	Aliases: []string{"repl", "rpl"},
	Desc:    "simple and quickly replace text or file contents",
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&trOpts)
		c.AddArg("text", "input text contents for process. allow @c,@FILE").WithAfterFn(func(a *gflag.CliArg) error {
			trOpts.text = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(trOpts.text)
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
			ret = strings.ReplaceAll(src, apputil.ResolveSep(trOpts.From), apputil.ResolveSep(trOpts.To))
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

// NewStringJoinCmd new command
func NewStringJoinCmd() *gcli.Command {
	var opt = struct {
		sep string
	}{}

	return &gcli.Command{
		Name:    "join",
		Aliases: []string{"j"},
		Desc:    "quick join multi line string by separator",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&opt.sep, "sep", "the separator for join, allow: NL,TAB", gflag.WithDefault(""))
			c.AddArg("text", "input strings for join", true, true)
		},
		Func: func(c *gcli.Command, _ []string) error {
			sep := apputil.ResolveSep(opt.sep)
			texts := c.Arg("text").Strings()
			for i, str := range texts {
				texts[i] = strings.ReplaceAll(str, "\n", sep)
			}

			fmt.Println(strings.Join(texts, sep))
			return nil
		},
	}
}
