package cmd

import (
	"github.com/gookit/gcli/v2"
	"github.com/gookit/gcli/v2/builtin"
	"github.com/gookit/kite/cmd/mkdown"
	"github.com/gookit/kite/cmd/swagger"
)

func AddCommands(app *gcli.App) {
	app.AddCommand(swagger.GenCode)
	app.AddCommand(swagger.DocBrowse)
	app.AddCommand(swagger.DocGen)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(builtin.GenAutoComplete())
	app.Add(mkdown.ConvertMD2html())
}
