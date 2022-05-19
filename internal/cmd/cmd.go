package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/internal/cmd/codegen"
	"github.com/inherelab/kite/internal/cmd/comtool"
	"github.com/inherelab/kite/internal/cmd/doctool"
	"github.com/inherelab/kite/internal/cmd/github"
	"github.com/inherelab/kite/internal/cmd/gitlab"
	"github.com/inherelab/kite/internal/cmd/gitx"
	"github.com/inherelab/kite/internal/cmd/gotool"
	"github.com/inherelab/kite/internal/cmd/mkdown"
	"github.com/inherelab/kite/internal/cmd/phptool"
	"github.com/inherelab/kite/internal/cmd/self"
	"github.com/inherelab/kite/internal/cmd/sql"
	"github.com/inherelab/kite/internal/cmd/swagger"
	"github.com/inherelab/kite/internal/cmd/taskx"
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
		taskx.TaskManage,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(pacutil.PacTools.WithHidden())
	app.Add(
		codegen.CodeGen,
		comtool.FileCat,
		comtool.FileFinder,
		comtool.BatchRun,
		comtool.HttpServe,
		comtool.RunScripts,
		comtool.HotReloadServe,
		builtin.GenAutoComplete(),
	)

	app.AddAliases("self:init", "init")
	app.AddAliases("self:info", "info")
}

func addListener(app *gcli.App) {
	app.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {

		// TODO
		return false
	})
}
