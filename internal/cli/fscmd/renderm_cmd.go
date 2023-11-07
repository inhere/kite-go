package fscmd

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil/finder"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/pkg/pkgutil"
)

// MultiRenderOpt options
type MultiRenderOpt struct {
	engine  string
	VarFmt  string         `flag:"desc=custom sets the variable format in template;shorts=vf;default={{,}}"`
	VarFile string         `flag:"desc=custom sets the variables file path. eg: --var-file tpl-vars.json"`
	Vars    gflag.KVString `flag:"desc=set template vars. allow multi, like: --vars name=Tom -v age=18;shorts=v"`

	// dir for template files
	Dir   string       `flag:"desc=the directory for find and render template files;shorts=d"`
	Exts  string       `flag:"desc=want render template files exts. multi by comma, like: .go,.md;shorts=ext"`
	Files gflag.String `flag:"desc=the template files. multi by comma, like: file1.tpl,file2.tpl;shorts=f"`

	// Include and Exclude match template files.
	// eg: --include name:*.go --exclude name:*_test.go
	Include gcli.Strings `flag:"desc=the include files rules;shorts=i,add"`
	Exclude gcli.Strings `flag:"desc=the exclude files rules;shorts=e,not"`

	Write bool `flag:"desc=write result to src file;shorts=w"`

	init bool
	vars *config.Config
}

// RenderFile render a template file
func (o *MultiRenderOpt) initVars() error {
	if o.init {
		return nil
	}

	o.init = true
	o.vars = pkgutil.NewConfig()

	// load vars from file
	if o.VarFile != "" {
		err := o.vars.LoadFiles(apputil.ResolvePath(o.VarFile))
		if err != nil {
			return err
		}
	}

	// load vars from cli
	o.vars.LoadSMap(o.Vars.SMap)

	return nil
}

// RenderFile render a template file
func (o *MultiRenderOpt) RenderFile(fPath string) error {
	if err := o.initVars(); err != nil {
		return err
	}

	// TODO
	return nil
}

// NewRenderMultiCmd create a command
func NewRenderMultiCmd() *gcli.Command {
	var opts = MultiRenderOpt{
		Vars: cflag.NewKVString(),
	}

	return &gcli.Command{
		Name:    "render-multi",
		Desc:    "render multi template files at once, allow use glob pattern, directory path",
		Aliases: []string{"renders", "mtpl"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&opts)

			c.StrOpt2(&opts.engine, "engine, eng", `select the template engine for rendering contents.
<b>Allow</>:
  go/go-tpl         - will use go template engine, support expression and control flow
  lite/lite-tpl     - will use lite template, support pipe expression, but not support control flow
  simple/replace    - only support simple variables replace rendering
`, gflag.WithDefault("lite"))

		},
		Func: func(c *gcli.Command, args []string) error {
			dump.P(opts)

			if opts.Files != "" {
				for _, fPath := range opts.Files.Strings() {
					fmt.Println(fPath)
				}
			}

			// find in dir
			if opts.Dir != "" {
				ff := finder.NewFinder(opts.Dir).
					WithExts(strutil.Split(opts.Exts, ","))
				// set finder options
				ff.IncludeRules(opts.Include.Strings())
				ff.ExcludeRules(opts.Exclude.Strings())

				for el := range ff.Find() {
					fmt.Println(el)
				}
			}

			return nil
		},
	}
}
