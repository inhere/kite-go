package cli

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/internal/cli/bintool"
	"github.com/inherelab/kite/internal/cli/codegen"
	"github.com/inherelab/kite/internal/cli/comtool"
	"github.com/inherelab/kite/internal/cli/doctool"
	"github.com/inherelab/kite/internal/cli/github"
	"github.com/inherelab/kite/internal/cli/gitlab"
	"github.com/inherelab/kite/internal/cli/gitx"
	"github.com/inherelab/kite/internal/cli/gotool"
	"github.com/inherelab/kite/internal/cli/mkdown"
	"github.com/inherelab/kite/internal/cli/phptool"
	"github.com/inherelab/kite/internal/cli/self"
	"github.com/inherelab/kite/internal/cli/sql"
	"github.com/inherelab/kite/internal/cli/swagger"
	"github.com/inherelab/kite/internal/cli/taskx"
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
		gitlab.GitLab,
		github.CmdForGithub,
		sql.SQLCmd,
		swagger.SwaggerCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
		phptool.PhpToolsCmd,
		self.KiteManage,
		self.KiteConf,
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

func addListener(app *gcli.App) {
	app.On(gcli.EvtCmdNotFound, func(ctx *gcli.HookCtx) bool {

		// TODO
		return false
	})
}
