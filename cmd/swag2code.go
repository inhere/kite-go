package cmd

import (
	"errors"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/ini/v2"
)

var swag2codeOpts = struct {
	tplDir     string
	SwagFile     string
	Template     string
	OutDir       string
	GroupSuffix  string
	ActionSuffix string
}{
	tplDir: "resource/templates",
}

var tplSuffix = ".tpl"

var tplFuncs = template.FuncMap{
	"join": strings.Join,
}

var Swag2code = &gcli.Command{
	Name:   "swag2code",
	UseFor: "generate rux API service codes by swagger.yaml or swagger.json",
	Config: func(c *gcli.Command) {
		c.StrOpt(&swag2codeOpts.SwagFile, "swagger-file", "f", "./swagger.json", "the swagger doc filepath")
		c.StrOpt(&swag2codeOpts.OutDir, "output", "o", "./gocodes", `the output directory for generated codes
if input 'stdout' will print codes on terminal
`)
		c.StrVar(&swag2codeOpts.Template, gcli.FlagMeta{
			Name:   "template",
			Desc:   "the template name for generate codes",
			Shorts: []string{"t"},
			DefVal: "rux-controller",
		})
		c.StrVar(&swag2codeOpts.GroupSuffix, gcli.FlagMeta{
			Name:   "group-suffix",
			Desc:   "Add suffix for group name. eg: API, Controller",
			DefVal: "API",
		})
		c.StrVar(&swag2codeOpts.ActionSuffix, gcli.FlagMeta{
			Name: "action-suffix",
			Desc: "Add suffix for action name. eg: Action, Method",
		})
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		swagFile := swag2codeOpts.SwagFile

		var bts []byte
		if swag.YAMLMatcher(swagFile) {
			bts, err = swag.YAMLDoc(swagFile)
		} else { // JSON file
			bts, err = ioutil.ReadFile(swagFile)
		}

		if err != nil {
			return err
		}

		doc := new(spec.Swagger)
		err = doc.UnmarshalJSON(bts)

		if len(doc.SwaggerProps.Paths.Paths) == 0 {
			return errors.New("API doc 'paths' is empty")
		}
		dump.Config(func(d *dump.Dumper) {
			d.MaxDepth = 8
		})

		tplDir := ini.String("swag2code.templateDir")
		if tplDir != "" {
			swag2codeOpts.tplDir = tplDir
		}

		generateByPathItem("/anything", doc.SwaggerProps.Paths.Paths["/anything"])

		return
	},
}

func generateByPathItem(path string, pathItem spec.PathItem) {
	tmpPath := strings.Trim(path, "/")
	group := tmpPath

	var substr string
	if strings.ContainsRune(tmpPath, '/') {
		nodes := strings.SplitN(tmpPath, "/", 2)
		group = nodes[0]
		substr = nodes[1]
	}

	dump.P(group, substr)
	dump.P(pathItem)

	tplFile := swag2codeOpts.Template
	if !strings.HasSuffix(tplFile, tplSuffix) {
		tplFile += tplSuffix
	}

	if !fsutil.IsFile(tplFile) {
		tplFile = filepath.Join(swag2codeOpts.tplDir, tplFile)
	}

	tpl := template.New("swag2code").Funcs(tplFuncs)
	template.Must(tpl.ParseFiles(tplFile))
}
