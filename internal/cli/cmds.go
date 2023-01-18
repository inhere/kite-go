package cli

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/inhere/kite/internal/cli/codegen"
	"github.com/inhere/kite/internal/cli/devcmd"
	"github.com/inhere/kite/internal/cli/doctool"
	"github.com/inhere/kite/internal/cli/fscmd"
	"github.com/inhere/kite/internal/cli/ghubcmd"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/internal/cli/glabcmd"
	"github.com/inhere/kite/internal/cli/gotool"
	"github.com/inhere/kite/internal/cli/mdcmd"
	"github.com/inhere/kite/internal/cli/phpcmd"
	"github.com/inhere/kite/internal/cli/pkgmanage"
	"github.com/inhere/kite/internal/cli/self"
	"github.com/inhere/kite/internal/cli/sqlcmd"
	"github.com/inhere/kite/internal/cli/strcmd"
	"github.com/inhere/kite/internal/cli/taskx"
	"github.com/inhere/kite/internal/cli/toolcmd"
	"github.com/inhere/kite/pkg/pacutil"
)

// Boot commands to gcli.App
func Boot(app *gcli.App) {
	addListener(app)

	Register(app)
}

// Register commands to gcli.App
func Register(app *gcli.App) {
	app.Add(
		doctool.DocumentCmd,
		gitcmd.GitCommands,
		glabcmd.GitLabCmd,
		ghubcmd.CmdForGithub,
		sqlcmd.SQLCmd,
		mdcmd.MkDownCmd,
		gotool.GoToolsCmd,
		phpcmd.PhpToolsCmd,
		strcmd.StringCmd,
		self.KiteManage,
		taskx.TaskManage,
		pkgmanage.ManageCmd,
		codegen.CodeGen,
		fscmd.FsCmd,
		toolcmd.ToolsCmd,
		toolcmd.RunScripts,
		devcmd.DevToolsCmd,
		builtin.GenAutoComplete(),
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(pacutil.PacTools.WithHidden())

	app.AddAliases("app:init", "init")
	app.AddAliases("app:info", "info")
	app.AddAliases("app:config", "conf", "config")
}

func addListener(app *gcli.App) {
	app.On(gcli.EvtCmdNotFound, func(ctx *gcli.HookCtx) bool {

		// TODO
		return false
	})
}
