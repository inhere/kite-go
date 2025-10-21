package fscmd

import (
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kautorw"
	"github.com/inhere/kite-go/pkg/util/bizutil"
)

type RenderFn func(src string, vars map[string]any) string

type templateCmdOpt struct {
	vars  gflag.KVString
	files string

	write  bool
	clipb  bool
	engine string
	varFmt string
	config string
	output string
	tplDir string
}

func (o *templateCmdOpt) loadConfig(varBox *config.Config) error {
	// load config file
	if o.config == "" {
		return nil
	}

	// varBox := bizutil.NewConfig()
	runConf := bizutil.NewConfig()
	cfgFile := app.PathMap.Resolve(o.config)
	err := runConf.LoadFiles(cfgFile)
	if err != nil {
		return err
	}

	// load custom vars from config file
	err = varBox.LoadData(runConf.Sub("vars"))
	if err != nil {
		return err
	}

	cfgSet := runConf.SubDataMap("settings")
	if !cfgSet.IsEmpty() {
		if v := cfgSet.Str("var_fmt"); len(v) > 0 {
			o.varFmt = v
		}
		if v := cfgSet.Str("engine"); len(v) > 0 {
			o.engine = v
		}
		if v := cfgSet.Str("tpl_dir"); len(v) > 0 {
			if strings.TrimLeft(v, "./") == "" {
				v = fsutil.DirPath(cfgFile)
			}
			o.tplDir = v
		}
	}
	show.AList("Loaded settings:", cfgSet)

	// load boot vars from input. eg: boot_var: [env, type]
	for _, name := range cfgSet.Strings("boot_var") {
		o.vars.IfValid(name, func(val string) {
			_ = varBox.Set(name, val)
		})
	}

	// load condition vars from input+config. eg: cond_var: [env, type]
	names := cfgSet.Strings("cond_var")
	for _, name := range names {
		o.vars.IfValid(name, func(val string) {
			_ = varBox.Set(name, val)
		})
		val := varBox.String(name)
		if val == "" {
			continue
		}

		// config key eg: env.qa
		err := varBox.LoadData(runConf.Sub(name + "." + val))
		if err != nil {
			return err
		}
	}

	// load united vars. eg: united_var: [[type, env], ...]
	ssList, _ := cfgSet.Slice("united_var")
	for _, item := range ssList {
		var nodes []string
		for _, name := range arrutil.MustToStrings(item) {
			if val := varBox.String(name); val != "" {
				nodes = append(nodes, val)
			} else {
				nodes = nil
				break
			}
		}

		if len(nodes) > 0 {
			// config key: name.{name}-{env}. eg: type=php,env=qa key=type.php-qa
			err := varBox.LoadData(runConf.Sub(strings.Join(nodes, ".")))
			if err != nil {
				return err
			}
		}
	}

	// load custom vars from command line
	if !o.vars.IsEmpty() {
		varBox.LoadSMap(o.vars.Data())
	}
	return nil
}

func (o *templateCmdOpt) makeEng() (bizutil.RenderFn, error) {
	return bizutil.NewTxtRender(o.engine, o.varFmt)
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

			c.StrOpt2(&ttOpts.files, "files,tpl", "set template file(s) for rendering, multiple use ',' to split")
			c.AddArg("files", "set template file(s) for rendering. same of --files", false, true)
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
		Func: func(c *gcli.Command, _ []string) error {
			varBox := bizutil.NewConfig()
			// load config file
			err := ttOpts.loadConfig(varBox)
			if err != nil {
				return err
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
			if ttOpts.files != "" {
				tplFiles = append(tplFiles, strings.Split(ttOpts.files, ",")...)
			}
			if len(tplFiles) == 0 {
				return errorx.E("please input template file(s)")
			}

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
