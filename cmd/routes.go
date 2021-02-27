package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kit/cmd/mkdown"
	"github.com/inherelab/kit/cmd/swagger"
)

func AddCommands(app *gcli.App) {
	app.Add(
		swagger.SwagCommand,
	)

	app.Add(
		mkdown.MkDownCmd,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(builtin.GenAutoComplete())
}
