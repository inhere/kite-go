package main

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/cmd/swagx/swagcmd"
)

// https://github.com/go-openapi
// https://github.com/getkin/kin-openapi
// https://github.com/pb33f/libopenapi
func main() {
	SwaggerCmd.MustRun(nil)
}

// SwaggerCmd instance
var SwaggerCmd = &gcli.Command{
	Name: "swag",
	Desc: "provide tools for use swagger/openapi",
	Subs: []*gcli.Command{
		swagcmd.DocBrowse,
		swagcmd.DocGen,
		swagcmd.Doc2MkDown,
		swagcmd.GenCode,
		swagcmd.InstallSwagGo,
	},
}
