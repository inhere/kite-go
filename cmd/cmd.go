package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/cmd/codegen"
	"github.com/inherelab/kite/cmd/comtool"
	"github.com/inherelab/kite/cmd/doctool"
	"github.com/inherelab/kite/cmd/github"
	"github.com/inherelab/kite/cmd/gitlab"
	"github.com/inherelab/kite/cmd/gitx"
	"github.com/inherelab/kite/cmd/gotool"
	"github.com/inherelab/kite/cmd/mkdown"
	"github.com/inherelab/kite/cmd/phptool"
	"github.com/inherelab/kite/cmd/self"
	"github.com/inherelab/kite/cmd/sql"
	"github.com/inherelab/kite/cmd/swagger"
	"github.com/inherelab/kite/pkg/pacutil"
)

// Boot commands to gcli.App
func Boot(app *gcli.App) {
	addListener(app)

	Register(app)
}

// Register commands to gcli.App
func Register(app *gcli.App) {
	// app.Add(
	// 	self.InitKite,
	// )

	app.Add(
		doctool.DocumentCmd,
		gitx.GitCommands,
		gitx.GitFlow,
		gitlab.CmdForGitlab,
		github.CmdForGithub,
		sql.SQLCmd,
		swagger.SwaggerCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
		phptool.PhpToolsCmd,
		self.KiteManage,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(pacutil.PacTools.WithHidden())
	app.Add(
		codegen.CodeGen,
		comtool.RunScripts,
		comtool.HotReloadServe,
		builtin.GenAutoComplete(),
	)

	app.AddAliases("self:init", "init")
}

func addListener(app *gcli.App) {
	app.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {

		// TODO
		return false
	})
}
