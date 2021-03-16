package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/cmd/comtool"
	"github.com/inherelab/kite/cmd/doctool"
	"github.com/inherelab/kite/cmd/git"
	"github.com/inherelab/kite/cmd/gotool"
	"github.com/inherelab/kite/cmd/mkdown"
	"github.com/inherelab/kite/cmd/sql"
	"github.com/inherelab/kite/cmd/swagger"
)

// Register commands to gcli.App
func Register(app *gcli.App) {
	app.Add(
		doctool.DocumentCmd,
		git.CmdsOfGit,
		sql.SQLCmd,
		swagger.SwaggerCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(
		comtool.HotReloadServe,
		builtin.GenAutoComplete(),
	)
}
