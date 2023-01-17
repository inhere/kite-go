package cli

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inhere/kite/internal/cli/codegen"
	"github.com/inhere/kite/internal/cli/comtool"
	"github.com/inhere/kite/internal/cli/devcmd"
	"github.com/inhere/kite/internal/cli/doctool"
	"github.com/inhere/kite/internal/cli/fscmd"
	"github.com/inhere/kite/internal/cli/github"
	"github.com/inhere/kite/internal/cli/gitlab"
	"github.com/inhere/kite/internal/cli/gitx"
	"github.com/inhere/kite/internal/cli/gotool"
	"github.com/inhere/kite/internal/cli/mkdown"
	"github.com/inhere/kite/internal/cli/phptool"
	"github.com/inhere/kite/internal/cli/pkgmanage"
	"github.com/inhere/kite/internal/cli/self"
	"github.com/inhere/kite/internal/cli/sql"
	"github.com/inhere/kite/internal/cli/taskx"
	"github.com/inhere/kite/pkg/pacutil"
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
		gitlab.GitLab,
		github.CmdForGithub,
		sql.SQLCmd,
		mkdown.MkDownCmd,
		gotool.GoToolsCmd,
		phptool.PhpToolsCmd,
		self.KiteManage,
		self.KiteConf,
		taskx.TaskManage,
		pkgmanage.ManageCmd,
		codegen.CodeGen,
		fscmd.FsCmd,
		comtool.ToolsCmd,
		comtool.BatchRun,
		comtool.RunScripts,
		devcmd.DevToolsCmd,
		builtin.GenAutoComplete(),
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(pacutil.PacTools.WithHidden())

	app.AddAliases("self:init", "init")
	app.AddAliases("self:info", "info")
}

func addListener(app *gcli.App) {
	app.On(gcli.EvtCmdNotFound, func(ctx *gcli.HookCtx) bool {

		// TODO
		return false
	})
}
