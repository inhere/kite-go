package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/gookit/color"
	"github.com/gookit/goutil/fsutil"
)

// SwagDoc struct
type SwagDoc struct {
	spec.Swagger
}

// PrintNode JSON string
func (d SwagDoc) PrintNode(name, filter string) (err error) {
	name = strings.ToLower(name)

	filtered := filter != ""
	showName := name
	indent := "  "

	var bts []byte
	switch name {
	case "tag", "tags":
		bts, err = json.MarshalIndent(d.Tags, "", indent)
	case "info":
		bts, err = json.MarshalIndent(d.Info, "", indent)
	case "path", "paths":
		showName = "paths"
		bts, err = json.MarshalIndent(d.Paths, "", indent)
	case "path-name", "path-names", "pathname", "pathnames":
		showName = "pathNames"

		var pathNames []string
		for path := range d.Paths.Paths {
			if filtered && strings.Contains(path, filter){
				pathNames = append(pathNames, path)
			}
		}
		if len(pathNames) > 0 {
			sort.Strings(pathNames)
		}

		bts, err = json.MarshalIndent(pathNames, "", indent)
	case "def", "defs", "definitions":
		showName = "definitions"
		bts, err = json.MarshalIndent(d.Definitions, "", indent)
	case "param", "params", "parameters":
		showName = "parameters"
		bts, err = json.MarshalIndent(d.Parameters, "", indent)
	case "res", "resp", "responses":
		showName = "responses"
		bts, err = json.MarshalIndent(d.Responses, "", indent)
	default:
		err = errors.New("node name value is invalid")
	}

	suffix := ""
	if err == nil {
		if filtered {
			suffix = "(filtered)"
		}

		color.Success.Printf("'%s' of the Document%s:\n", showName, suffix)
		fmt.Println(string(bts))
	}

	return
}

func (d SwagDoc) GetTagInfo(name string) (spec.Tag, bool) {
	for _, tag := range d.Tags {
		if name != tag.Name {
			return tag, true
		}
	}
	return spec.Tag{}, false
}

var swagger = &SwagDoc{}

func loadDocFile(swagFile string) (err error) {
	if !fsutil.IsFile(swagFile) {
		return errors.New("the swagger file not exist")
	}

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
	err = swagger.UnmarshalJSON(bts)
	return
}
