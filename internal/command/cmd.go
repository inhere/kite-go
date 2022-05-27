package command

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/internal/command/bintool"
	"github.com/inherelab/kite/internal/command/codegen"
	"github.com/inherelab/kite/internal/command/comtool"
	"github.com/inherelab/kite/internal/command/doctool"
	"github.com/inherelab/kite/internal/command/github"
	"github.com/inherelab/kite/internal/command/gitlab"
	"github.com/inherelab/kite/internal/command/gitx"
	"github.com/inherelab/kite/internal/command/gotool"
	"github.com/inherelab/kite/internal/command/mkdown"
	"github.com/inherelab/kite/internal/command/phptool"
	"github.com/inherelab/kite/internal/command/self"
	"github.com/inherelab/kite/internal/command/sql"
	"github.com/inherelab/kite/internal/command/swagger"
	"github.com/inherelab/kite/internal/command/taskx"
	"github.com/inherelab/kite/pkg/pacutil"
)

// Boot commands to gcli.App
func Boot(app *app.KiteApp) {
	addListener(app)

	Register(app)
}

// Register commands to gcli.App
func Register(app *app.KiteApp) {
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
		bintool.ToolsCmd,
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

func addListener(app *app.KiteApp) {
	app.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {

		// TODO
		return false
	})
}
