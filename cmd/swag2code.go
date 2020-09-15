package cmd

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/ini/v2"
	"github.com/gookit/slog"
)

const stdoutName = "stdout"
// refer links:
// https://petstore.swagger.io/
// https://petstore3.swagger.io/
var swag2codeOpts = struct {
	tplDir string
	writer io.Writer
	// options
	SwagFile     string
	Template     string
	OutDir       string
	GroupSuffix  string
	ActionSuffix string
	SpecPaths    gcli.Strings
}{
	tplDir: "resource/templates",
}

var tplSuffix = ".tpl"
var tplEngine *template.Template

var tplFuncs = template.FuncMap{
	"join": strings.Join,
	"upper": strings.ToUpper,
	"upFirst": strutil.UpperFirst,
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
		c.VarOpt(&swag2codeOpts.SpecPaths, "paths", "", "only generate for the given spec paths")
		c.Examples = `
{$fullCmd} --paths /anything /blog/{id}
`
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		// load swagger doc file
		doc, err := loadSwaggerDocFile()
		if err != nil {
			return err
		}

		// check paths
		if len(doc.SwaggerProps.Paths.Paths) == 0 {
			return errors.New("API doc 'paths' is empty")
		}

		// load and parse template
		if err = loadAndParseTemplate(); err != nil {
			return
		}

		dump.Config(func(d *dump.Dumper) {
			d.MaxDepth = 8
		})

		// only generate special paths.
		if len(swag2codeOpts.SpecPaths) > 0 {
			slog.Info("will only generate for the paths: ", swag2codeOpts.SpecPaths)

			for _, path := range swag2codeOpts.SpecPaths {
				if pathItem, ok := doc.SwaggerProps.Paths.Paths[path]; ok {
					err = generateByPathItem(path, pathItem)
					if err != nil {
						return
					}
				} else { // path not exists
					slog.Errorf("- the path '%s' is not exists on docs, skip gen", path)
				}
			}
		} else {
			for path, pathItem := range doc.SwaggerProps.Paths.Paths {
				err = generateByPathItem(path, pathItem)
				if err != nil {
					return
				}
			}
		}
		return
	},
}

// ActionItem struct
type ActionItem struct {
	// Path route path
	Path string
	Tags []string
	// METHOD the request METHOD. eg. GET
	METHOD string
	// MethodName the action method name
	MethodName string
	// MethodDesc the action method desc
	MethodDesc string
}

func loadSwaggerDocFile() (doc spec.Swagger, err error) {
	swagFile := swag2codeOpts.SwagFile

	var bts []byte
	if swag.YAMLMatcher(swagFile) {
		bts, err = swag.YAMLDoc(swagFile)
	} else { // JSON file
		bts, err = ioutil.ReadFile(swagFile)
	}

	if err != nil {
		return
	}

	// doc := new(spec.Swagger)
	err = doc.UnmarshalJSON(bts)
	return
}

func loadAndParseTemplate() error  {
	tplDir := ini.String("swag2code.templateDir")
	if tplDir != "" {
		swag2codeOpts.tplDir = tplDir
	}

	tplFile := swag2codeOpts.Template
	if !strings.HasSuffix(tplFile, tplSuffix) {
		tplFile += tplSuffix
	}

	if !fsutil.IsFile(tplFile) {
		// find from default template dir
		tplFile = filepath.Join(swag2codeOpts.tplDir, tplFile)
		if !fsutil.IsFile(tplFile) {
			return errors.New("template file not exists, file: " + swag2codeOpts.Template)
		}
	}

	slog.Infof("use the template file: %s", tplFile)

	bts, err := ioutil.ReadFile(tplFile)
	if err != nil {
		return err
	}
	root := template.New("swag2code")

	slog.Info("do parsing the template contents")

	// tpl, err := tplEngine.New(swag2codeOpts.Template).Funcs(tplFuncs).ParseFiles(tplFile)
	tplEngine, err = root.New(swag2codeOpts.Template).Funcs(tplFuncs).Parse(string(bts))

	return err
}

func generateByPathItem(path string, pathItem spec.PathItem) (err error) {
	var subpath string
	tmpPath := strings.Trim(path, "/")

	group := tmpPath
	if strings.ContainsRune(tmpPath, '/') {
		nodes := strings.SplitN(tmpPath, "/", 2)

		group, subpath = nodes[0], nodes[1]
	}

	// dump.P(group, subpath)
	// dump.P(pathItem)
	slog.Infof("will generate for the path '%s'. (group: %s, subpath: %s)", path, group, subpath)

	data := &struct {
		GroupName string
		GroupDesc string
		GroupPath string
		Actions   []ActionItem
	} {
		GroupName: group,
		GroupDesc: "",
		GroupPath: path,
	}

	piProps := pathItem.PathItemProps
	if piProps.Get != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodGet, piProps.Get))
	}
	if piProps.Put != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodPut, piProps.Put))
	}
	if piProps.Head != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodPut, piProps.Head))
	}
	if piProps.Post != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodPost, piProps.Post))
	}
	if piProps.Patch != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodPatch, piProps.Patch))
	}
	if piProps.Delete != nil {
		data.Actions = append(data.Actions, buildActionItem(http.MethodDelete, piProps.Delete))
	}

	// dump.P(pathItem.PathItemProps)

	slog.Info("do rendering the template contents")

	buf := new(bytes.Buffer)
	err = tplEngine.Execute(buf, data)
	if err == nil {
		slog.Infof("output the generated contents to: %s", swag2codeOpts.OutDir)

		// dump to stdout
		if swag2codeOpts.OutDir == stdoutName {
			_, err = buf.WriteTo(os.Stdout)
		}
	}

	return
}

func buildActionItem(method string, op *spec.Operation) ActionItem {
	return ActionItem{
		Path:       "",
		METHOD:     method,
		MethodName: strutil.UpperFirst(strings.ToLower(method)),
		MethodDesc: op.Description,
	}
}
