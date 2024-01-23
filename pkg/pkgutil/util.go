package pkgutil

import (
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
)

// NewConfig box instance
func NewConfig() *config.Config {
	return config.
		NewWithOptions("kite", config.ParseEnv, config.ParseDefault, config.WithTagName("json")).
		WithDriver(yaml.Driver, JSON5Driver, ini.Driver)
}

// RenderFn render function
type RenderFn func(src string, vars map[string]any) string

// NewTxtRender create a text render function
func NewTxtRender(engine, varFmt string) (RenderFn, error) {
	switch engine {
	case "jet": // github.com/CloudyKit/jet/v6
		left, right := strutil.TrimCut(varFmt, ",")
		jte := jet.NewSet(jet.NewInMemLoader(), jet.WithDelims(left, right))

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
			opt.SetVarFmt(varFmt)
		})
		return tplE.RenderString, nil
	case "simple", "replace":
		return textutil.NewVarReplacer(varFmt).Replace, nil
	default:
		return nil, errorx.Errf("invalid engine name %q", engine)
	}
}
