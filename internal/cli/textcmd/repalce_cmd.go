package textcmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/pkg/kautorw"
	"github.com/inhere/kite-go/pkg/pkgutil"
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

		var dst string
		if trOpts.Regex {
			reg := regexp.MustCompile(trOpts.From)
			dst = reg.ReplaceAllString(src, trOpts.To)
		} else {
			dst = strings.ReplaceAll(src, trOpts.From, trOpts.To)
		}

		fmt.Println(dst)
		return nil
	},
}

var ttOpts = struct {
	vars gflag.KVString
	text string

	engine  string
	varFmt  string
	varFile string
	output  string
}{
	engine: "simple",
	vars:   cflag.NewKVString(),
}

// NewTemplateCmd instance
func NewTemplateCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "render",
		Aliases: []string{"tpl-render"},
		Desc:    "simple rendering text template contents by replace",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&ttOpts.varFmt, "var-fmt", "custom sets the variable format in template", gflag.WithDefault("{{,}}"))
			c.StrOpt2(&ttOpts.varFile, "var-file", "custom sets the variables file path")
			c.StrOpt2(&ttOpts.output, "output,o", "custom sets the output target", gflag.WithDefault("stdout"))
			c.VarOpt2(&ttOpts.vars, "vars,var,v", "sets template variables for render. format: `KEY=VALUE`")

			c.StrOpt2(&ttOpts.engine, "engine, eng", `select the template engine for rendering contents. 
<b>Allow</>:
  go/go-tpl         - will use go template engine and support expression
  simple/replace    - only support simple variables replace rendering
`)

			c.AddArg("text", "template file or contents for rendering").WithAfterFn(func(a *gflag.CliArg) error {
				ttOpts.text = a.String()
				return nil
			})
		},
		Help: `
## simple example
  {$fullCmd} -v name=inhere -v age=234 'hi, {{name}}, age is {{ age }}'

## go-tpl example
  {$fullCmd} --eng go-tpl -v name=inhere -v age=234 'hi, {{.name}}, age is {{ .age }}'

## use template file
  {$fullCmd} --var-file /path/to/_variables.yaml @/path/to/my-template.tpl
`,
		Func: func(c *gcli.Command, _ []string) error {
			src, err := apputil.ReadSource(ttOpts.text)
			if err != nil {
				return err
			}

			varBox := pkgutil.NewConfig()
			if ttOpts.varFile != "" {
				err := varBox.LoadFiles(ttOpts.varFile)
				if err != nil {
					return err
				}
			}

			if len(ttOpts.vars.Data()) > 0 {
				varBox.LoadSMap(ttOpts.vars.Data())
			}
			show.AList("Loaded variables:", varBox.Data())

			// do rendering
			ret := src
			switch ttOpts.engine {
			case "go", "go-tpl":
				ret = strutil.RenderTemplate(src, varBox.Data(), nil)
			case "simple", "replace":
				ret = textutil.ReplaceVars(src, varBox.Data(), ttOpts.varFmt)
			default:
				return c.NewErrf("invalid engine name %q", ttOpts.engine)
			}

			sw := kautorw.NewSourceWriter(ttOpts.output)
			return sw.WriteString(ret)
		},
	}
}
