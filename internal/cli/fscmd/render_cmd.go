package fscmd

import (
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kautorw"
	"github.com/inhere/kite-go/pkg/pkgutil"
)

type RenderFn func(src string, vars map[string]any) string

type templateCmdOpt struct {
	vars gflag.KVString
	// files []string

	write  bool
	clipb  bool
	engine string
	varFmt string
	config string
	output string
	tplDir string
}

var jte = jet.NewSet(jet.NewInMemLoader())

func (o *templateCmdOpt) makeEng() (RenderFn, error) {
	switch o.engine {
	case "jet": // github.com/CloudyKit/jet/v6
		return func(src string, vars map[string]any) string {
			jt, err := jte.Parse("temp-file.jet", src)
			if err != nil {
				return err.Error()
			}

			buf := new(strings.Builder)
			err = jt.Execute(buf, nil, vars)
			if err != nil {
				return err.Error()
			}
			return buf.String()
		}, nil
	case "go", "go-tpl":
		return func(src string, vars map[string]any) string {
			return textutil.RenderGoTpl(src, vars)
		}, nil
	case "lite", "lite-tpl":
		tplE := textutil.NewLiteTemplate(func(opt *textutil.LiteTemplateOpt) {
			opt.SetVarFmt(o.varFmt)
		})
		return tplE.RenderString, nil
	case "simple", "replace":
		return textutil.NewVarReplacer(o.varFmt).Replace, nil
	default:
		return nil, errorx.Rawf("invalid engine name %q", o.engine)
	}
}

// NewTemplateCmd instance
func NewTemplateCmd() *gcli.Command {
	var ttOpts = templateCmdOpt{
		engine: "lite",
		vars:   cflag.NewKVString(),
	}

	return &gcli.Command{
		Name:    "render",
		Aliases: []string{"tpl", "tpl-render"},
		Desc:    "quickly rendering given template files and with variables",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&ttOpts.varFmt, "var-fmt", "custom sets the variable format in template", gflag.WithDefault("{{,}}"))
			c.StrOpt2(&ttOpts.config, "config,c", "custom config. allow sets the variables file path")
			c.StrOpt2(&ttOpts.output, "output,o", "custom sets the output target", gflag.WithDefault("stdout"))
			c.VarOpt2(&ttOpts.vars, "vars,var,v", "sets template variables for render. format: `KEY=VALUE`")
			c.BoolOpt2(&ttOpts.write, "write,w", "write result to src file, on input is filepath")
			c.BoolOpt2(&ttOpts.clipb, "clip,cb", "write result to the system clipboard")

			c.StrOpt2(&ttOpts.engine, "engine, eng", `select the template engine for rendering contents.
<b>Allow</>:
  jet               - use CloudyKit/jet template engine, support expression and control flow.
  go/go-tpl         - will use go template engine, support expression and control flow.
  lite/lite-tpl     - will use lite template, support pipe expression, but not support control flow
  simple/replace    - only support simple variables replace rendering
`)

			c.AddArg("files", "set template file(s) for rendering", true, true)
		},
		Help: `
## Note
 - support use path alias(kite,user). eg: @tpl_dir/some.tpl

## Output
 - default output to stdout. same of -o=@stdout
 - use --write option, will write result to src file, on input is filepath. same of -o=@src

 ### use expr:
   vars in expr:
    - $fileName   - the file name of current file.
    - $nameNoExt  - the file name of current file, not contains extension.

   Example:
    '@workdir/$fileName' will write result to workdir/$fileName

## simple example
  {$fullCmd} -v name=inhere -v age=234 'hi, {{name}}, age is {{ age }}'

## go-tpl example
  {$fullCmd} --eng go-tpl -v name=inhere -v age=234 'hi, {{.name}}, age is {{ .age }}'

## use template file
  {$fullCmd} --config /path/to/_config.yaml /path/to/my-template.tpl
`,
		Func: func(c *gcli.Command, _ []string) (err error) {
			varBox := pkgutil.NewConfig()
			runConf := pkgutil.NewConfig()
			// load config file
			if ttOpts.config != "" {
				cfgFile := app.PathMap.Resolve(ttOpts.config)
				err = runConf.LoadFiles(cfgFile)
				if err != nil {
					return err
				}

				// load custom vars from config file
				cfgVars := runConf.SubDataMap("vars")
				err = varBox.LoadData(map[string]any(cfgVars))
				if err != nil {
					return err
				}

				cfgSet := runConf.SubDataMap("settings")
				if cfgSet != nil {
					if v := cfgSet.Str("var_fmt"); len(v) > 0 {
						ttOpts.varFmt = v
					}
					if v := cfgSet.Str("engine"); len(v) > 0 {
						ttOpts.engine = v
					}
					if v := cfgSet.Str("tpl_dir"); len(v) > 0 {
						if strings.TrimLeft(v, "./") == "" {
							v = fsutil.DirPath(cfgFile)
						}
						ttOpts.tplDir = v
					}
				}
				show.AList("Loaded settings:", cfgSet)
			}

			if len(ttOpts.vars.Data()) > 0 {
				varBox.LoadSMap(ttOpts.vars.Data())
			}
			show.AList("Loaded variables:", varBox.Data(), func(opts *show.ListOption) {
				opts.IgnoreEmpty = false
			})

			if len(ttOpts.tplDir) > 0 {
				ttOpts.tplDir = app.PathMap.Resolve(ttOpts.tplDir)
				app.PathMap.AddAlias("tpl_dir", ttOpts.tplDir)
			}

			// do rendering
			engFn, err := ttOpts.makeEng()
			if err != nil {
				return err
			}

			sw := kautorw.NewSourceWriter(ttOpts.output)
			if ttOpts.write {
				sw.WithDst("@src")
			} else if ttOpts.clipb {
				sw.WithDst("@clip")
			}

			tplFiles := cliutil.SplitMulti(c.Arg("files").Strings(), ",")
			openFlush := len(tplFiles) > 1 && sw.DstType() == kautorw.TypeClip

			for _, file := range tplFiles {
				file = app.PathMap.Resolve(file)
				body, err := fsutil.ReadStringOrErr(file)
				if err != nil {
					c.Errorf("read file error: %s", err.Error())
					continue
				}

				// rendering contents
				str := engFn(body, varBox.Data())

				sw.SetSrcFile(file)
				if err = sw.WriteString(str); err != nil {
					return err
				}
			}

			if openFlush {
				return sw.StopFlush()
			}
			return nil
		},
	}
}
