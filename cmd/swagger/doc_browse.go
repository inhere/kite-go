package swagger

import (
	"github.com/gookit/gcli/v2"
)

var docBrowseOpts = struct {
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
	},
	Func: func(c *gcli.Command, args []string) (err error) {
		// load swagger doc file
		if err := loadDocFile(docBrowseOpts.SwagFile); err != nil {
			return err
		}

		return swagger.PrintNode(docBrowseOpts.NodeName)
	},
}
