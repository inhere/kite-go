package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kit/cmd/gotool"
	"github.com/inherelab/kit/cmd/mkdown"
	"github.com/inherelab/kit/cmd/swagger"
)

// Register commands to gcli.App
func Register(app *gcli.App) {
	app.Add(
		swagger.SwaggerCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(builtin.GenAutoComplete())
}
