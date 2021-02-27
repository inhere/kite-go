package swagger

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/gookit/gcli/v3"
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
	OutDir   string
	PkgName  string
	SwagFile string
	Template string
	// GroupType paths group type:
	// none - dont group
	// path - by path
	// tag - by tag name
	GroupType    string
	GroupSuffix  string
	ActionSuffix string
	SpecPaths    gcli.Strings
}{
	tplDir: "resource/template/codegen", // default dir
}

var (
	tplEngine *template.Template
	tplSuffix = ".tpl"
	tplFuncs  = template.FuncMap{
		"join":    strings.Join,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"upFirst": strutil.UpperFirst,
	}
)

var GenCode = &gcli.Command{
	Name:    "swag2go",
	Aliases: []string{"swag2code"},
	Desc:  "generate go API service codes by swagger.yaml or swagger.json",
	Config: func(c *gcli.Command) {
		c.StrOpt(&swag2codeOpts.SwagFile,
			"swagger-file",
			"f",
			"api",
			"the swagger doc filepath",
		)
		c.StrOpt(&swag2codeOpts.PkgName, "pgk-name", "p", "./swagger.json", "the generated package name")
		c.StrOpt(&swag2codeOpts.OutDir,
			"output",
			"o",
			"./gocodes",
			`the output directory for generated codes
if input 'stdout' will print codes on terminal
`)
		c.StrVar(&swag2codeOpts.Template, gcli.FlagMeta{
			Name:   "template",
			Desc:   "the template name for generate codes",
			Shorts: []string{"t"},
			DefVal: "rux-controller",
		})
		c.StrVar(&swag2codeOpts.GroupType, gcli.FlagMeta{
			Name: "group-type",
			Desc: `the code generate group type. allow:
none - dont generate group struct
path - group by path name
tag  - group by tag name
`,
			DefVal: "tag",
			Validator: func(val string) error {
				ss := []string{"none", "path", "tag"}
				for _, s := range ss {
					if s == val {
						return nil
					}
				}

				return fmt.Errorf("'group-type' must one of %v", ss)
			},
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
{$fullCmd} -f testdata/swagger.json --paths /anything/{anything} --group-type path -o stdout
`
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		// load swagger doc file
		if err := loadDocFile(swag2codeOpts.SwagFile); err != nil {
			return err
		}

		// check paths
		if len(swagger.SwaggerProps.Paths.Paths) == 0 {
			return errors.New("API doc 'paths' is empty")
		}

		// load and parse template
		if err = loadAndParseTemplate(); err != nil {
			return
		}

		dump.Config(func(d *dump.Dumper) {
			d.MaxDepth = 8
		})

		outDir := swag2codeOpts.OutDir
		if outDir != stdoutName && fsutil.IsDir(outDir) {
			slog.Info("create the output directory:", outDir)

			err = os.MkdirAll(outDir, 0664)
			if err != nil {
				return err
			}
		}

		// only generate special paths.
		if len(swag2codeOpts.SpecPaths) > 0 {
			slog.Info("will only generate for the paths: ", swag2codeOpts.SpecPaths)

			for _, path := range swag2codeOpts.SpecPaths {
				if pathItem, ok := swagger.SwaggerProps.Paths.Paths[path]; ok {
					err = generateByPathItem(path, pathItem)
					if err != nil {
						return
					}
				} else { // path not exists
					slog.Errorf("- the path '%s' is not exists on docs, skip gen", path)
				}
			}
		} else {
			for path, pathItem := range swagger.SwaggerProps.Paths.Paths {
				err = generateByPathItem(path, pathItem)
				if err != nil {
					return
				}
			}
		}
		return
	},
}

func loadAndParseTemplate() error {
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

type GroupData struct {
	GroupName string
	GroupDesc string
	GroupPath string

	PkgName string
	Actions []ActionItem
	TagMap  map[string]string
}

func (gd *GroupData) addActionItem(method string, op *spec.Operation) {
	desc := op.Description
	if desc == "" {
		if op.Summary != "" {
			desc = op.Summary
		} else {
			var prefix string
			switch method {
			case http.MethodGet:
				prefix = "query"
			case http.MethodPost:
				prefix = "create"
			case http.MethodPut, http.MethodPatch:
				prefix = "update"
			default:
				prefix = "do"
			}
			desc = prefix + " operation for the " + gd.GroupName
		}
	}

	// dump.P(op.OperationProps)
	actionName := strutil.UpperFirst(strings.ToLower(method))

	item := ActionItem{
		Path:   "",
		Tags:   op.Tags,
		METHOD: method,
		// other
		MethodName: actionName,
		MethodDesc: desc,
	}

	gd.Actions = append(gd.Actions, item)

	// load tag desc
	if gd.GroupDesc == "" && len(op.Tags) > 0 {
		tag, ok := swagger.GetTagInfo(op.Tags[0])
		if ok {
			gd.GroupDesc = tag.Description
		}
	}
}

func (gd *GroupData) appendTags(tags []string) {

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

// TagComments line
func (a ActionItem) TagComments() string {
	if len(a.Tags) == 0 {
		return ""
	}

	return fmt.Sprintf("\n// @Tags %s", a.Tags[0])
}

func generateByPathItem(path string, pathItem spec.PathItem) (err error) {
	var subpath string
	tmpPath := strings.Trim(path, "/")

	grpName := tmpPath
	if strings.ContainsRune(tmpPath, '/') {
		nodes := strings.SplitN(tmpPath, "/", 2)

		grpName, subpath = nodes[0], nodes[1]
	}

	// dump.P(group, subpath)
	// dump.P(pathItem)
	slog.Infof("will generate for the path '%s'. (group: %s, subpath: %s)", path, grpName, subpath)

	data := &GroupData{
		PkgName:   swag2codeOpts.PkgName,
		GroupName: grpName,
		GroupPath: path,
	}

	piProps := pathItem.PathItemProps
	if piProps.Get != nil {
		data.addActionItem(http.MethodGet, piProps.Get)
	}
	if piProps.Put != nil {
		data.addActionItem(http.MethodPut, piProps.Put)
	}
	if piProps.Head != nil {
		data.addActionItem(http.MethodPut, piProps.Head)
	}
	if piProps.Post != nil {
		data.addActionItem(http.MethodPost, piProps.Post)
	}
	if piProps.Patch != nil {
		data.addActionItem(http.MethodPatch, piProps.Patch)
	}
	if piProps.Delete != nil {
		data.addActionItem(http.MethodDelete, piProps.Delete)
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
