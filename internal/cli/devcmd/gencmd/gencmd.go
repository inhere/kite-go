package gencmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/ini/v2"
)

// CodeGenCmd instance
var CodeGenCmd = &gcli.Command{
	Name:    "gen",
	Desc:    "quickly code or project generate",
	Aliases: []string{"cgen", "generate"},
	Subs: []*gcli.Command{
		ParseTemplateCmd,
		NewProjectCmd(),
	},
}

var genOpts = struct {
	tpl string
}{}

// ParseTemplateCmd instance
var ParseTemplateCmd = &gcli.Command{
	Name: "parse",
	Desc: "parse template for generate code",
	Func: func(c *gcli.Command, args []string) error {
		if !fsutil.IsFile(genOpts.tpl) {
			return c.NewErr("template file not exists")
		}

		c.Infoln("will generate code by parse:", genOpts.tpl)

		buf := new(bytes.Buffer)
		bts := fsutil.MustReadFile(genOpts.tpl)

		nodes := strings.SplitN(string(bts), "\n###\n", 2)
		if len(nodes) != 2 {
			return c.NewErr("template content is invalid")
		}

		varCode, tplCode := nodes[0], nodes[1]

		c.Infoln("- parse template vars")
		err := ini.LoadStrings(varCode)
		if err != nil {
			return err
		}

		data := ini.StringMap("")
		// dump.Println(data)

		varMap := make(map[string]interface{}, len(data))
		for k, v := range data {
			// `[Info,Trace]` as array
			if v[0] == '[' {
				varMap[k] = strutil.Split(strings.Trim(v, "[]"), ",")
			} else {
				varMap[k] = v
			}
		}

		// open debug
		if gcli.IsDebugMode() {
			d := dump.NewWithOptions(func(opts *dump.Options) {
				opts.ShowFlag = dump.Fnopos
			})

			d.Println(varMap)
		}

		c.Infoln("- render template contents")
		t := template.New("parseTpl")
		template.Must(t.Parse(tplCode))

		err = t.Execute(buf, varMap)
		if err != nil {
			return err
		}

		fmt.Print(buf.String())
		return nil
	},

	Config: func(c *gcli.Command) {
		c.StrOpt(&genOpts.tpl, "tpl", "t", "", "the template file path")
	},
}
