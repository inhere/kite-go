package textcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/pkg/kautorw"
	"github.com/inhere/kite-go/pkg/pkgutil"
)

type templateCmdOpt struct {
	vars gflag.KVString
	text string

	write   bool
	engine  string
	varFmt  string
	varFile string
	output  string
}

// NewTemplateCmd instance
func NewTemplateCmd(mustFile bool) *gcli.Command {
	var ttOpts = templateCmdOpt{
		engine: "simple",
		vars:   cflag.NewKVString(),
	}

	return &gcli.Command{
		Name:    "render",
		Aliases: []string{"tpl-render"},
		Desc:    "simple rendering text template contents by replace",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&ttOpts.varFmt, "var-fmt", "custom sets the variable format in template", gflag.WithDefault("{{,}}"))
			c.StrOpt2(&ttOpts.varFile, "var-file", "custom sets the variables file path")
			c.StrOpt2(&ttOpts.output, "output,o", "custom sets the output target", gflag.WithDefault("stdout"))
			c.VarOpt2(&ttOpts.vars, "vars,var,v", "sets template variables for render. format: `KEY=VALUE`")
			c.BoolOpt2(&ttOpts.write, "write,w", "write result to src file, on input is filepath")

			c.StrOpt2(&ttOpts.engine, "engine, eng", `select the template engine for rendering contents.
<b>Allow</>:
  go/go-tpl         - will use go template engine, support expression and control flow
  lite/lite-tpl     - will use lite template, support pipe expression, but not support control flow
  simple/replace    - only support simple variables replace rendering
`)

			c.AddArg("text", "src template file or contents for rendering").WithAfterFn(func(a *gflag.CliArg) error {
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
			if mustFile && !fsutil.IsFile(ttOpts.text) {
				return c.NewErrf("the input is not a file: %s", ttOpts.text)
			}

			srr := apputil.NewSReader(ttOpts.text)
			if mustFile {
				srr.WithConfig(kautorw.WithDefaultAsFile())
			}

			src, err := srr.TryReadString()
			if err != nil {
				return err
			}

			varBox := pkgutil.NewConfig()
			// load config file
			if ttOpts.varFile != "" {
				err = varBox.LoadFiles(ttOpts.varFile)
				if err != nil {
					return err
				}
			}

			if len(ttOpts.vars.Data()) > 0 {
				varBox.LoadSMap(ttOpts.vars.Data())
			}
			show.AList("Loaded variables:", varBox.Data())

			// do rendering
			engFn, err := pkgutil.NewTxtRender(ttOpts.engine, ttOpts.varFmt)
			if err != nil {
				return err
			}
			ret := engFn(src, varBox.Data())

			sw := kautorw.NewSourceWriter(ttOpts.output)
			sw.SetSrcFile(ttOpts.text)

			if ttOpts.write {
				sw.WithDst("@src")
				if !sw.HasSrcFile() {
					return c.NewErrf("with option --write, but input is not a file")
				}
			}
			return sw.WriteString(ret)
		},
	}
}
