package swagger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
)

var docBrowseOpts = struct {
	Filter   string
	NodeName string
	PathName string
	SwagFile string
}{}

var DocBrowse = &gcli.Command{
	Name:    "swag:browse",
	Aliases: []string{"swag:cat", "swag:see"},
	UseFor:  "open browser for browse input swagger doc file",
	Config: func(c *gcli.Command) {
		c.StrOpt(&docBrowseOpts.SwagFile,
			"swagger-file", "f",
			"swagger.json",
			"the swagger document file path",
		)
		c.StrVar(&docBrowseOpts.NodeName, gcli.FlagMeta{
			Name:   "node",
			Shorts: []string{"n"},
			Desc:   "show parts of the the documents.\nallow: tags, info, paths, defs, responses",
			// must
			// Required: true,
		})
		c.StrVar(&docBrowseOpts.PathName, gcli.FlagMeta{
			Name:   "path",
			Shorts: []string{"p"},
			Desc:   "show path info of the the `documents.paths`. eg: /anything",
		})
		c.StrVar(&docBrowseOpts.Filter, gcli.FlagMeta{
			Name: "filter",
			Desc: "filter the results of path or node",
		})
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		// load swagger doc file
		if err := loadDocFile(docBrowseOpts.SwagFile); err != nil {
			return err
		}

		if docBrowseOpts.PathName != "" {
			path := docBrowseOpts.PathName
			if pItem, ok := swagger.SwaggerProps.Paths.Paths[path]; ok {
				bts, err := json.MarshalIndent(pItem, "", "  ")
				if err != nil {
					return err
				}

				color.Success.Printf("'paths.%s' of the Document:\n", path)
				fmt.Println(string(bts))
			} else {
				return fmt.Errorf("'paths.%s' is not exist of the Document", path)
			}
		}

		if docBrowseOpts.NodeName != "" {
			return swagger.PrintNode(docBrowseOpts.NodeName, docBrowseOpts.Filter)
		}

		return errors.New("please setting --node|--path value")
	},
}
