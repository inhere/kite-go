package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inherelab/kite/cmd/comtool"
	"github.com/inherelab/kite/cmd/doctool"
	"github.com/inherelab/kite/cmd/github"
	"github.com/inherelab/kite/cmd/gitlab"
	"github.com/inherelab/kite/cmd/gituse"
	"github.com/inherelab/kite/cmd/gotool"
	"github.com/inherelab/kite/cmd/mkdown"
	"github.com/inherelab/kite/cmd/self"
	"github.com/inherelab/kite/cmd/sql"
	"github.com/inherelab/kite/cmd/swagger"
)

// Register commands to gcli.App
func Register(app *gcli.App) {
	// app.Add(
	// 	self.InitKite,
	// )

	app.Add(
		doctool.DocumentCmd,
		gituse.GitCommands,
		gituse.GitFlow,
		gitlab.CmdForGitlab,
		github.CmdForGithub,
		sql.SQLCmd,
		swagger.SwaggerCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
		self.KiteManage,
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(
		comtool.RunScripts,
		comtool.HotReloadServe,
		builtin.GenAutoComplete(),
	)

	app.AddAliases("self:init", "init")
}
