package cli

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/gookit/gcli/v3/events"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/cli/codegen"
	"github.com/inhere/kite/internal/cli/devcmd"
	"github.com/inhere/kite/internal/cli/doctool"
	"github.com/inhere/kite/internal/cli/fscmd"
	"github.com/inhere/kite/internal/cli/ghubcmd"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/internal/cli/glabcmd"
	"github.com/inhere/kite/internal/cli/gocmd"
	"github.com/inhere/kite/internal/cli/httpcmd"
	"github.com/inhere/kite/internal/cli/javacmd"
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
		codegen.CodeGen,
		devcmd.DevToolsCmd,
		doctool.DocumentCmd,
		fscmd.FsCmd,
		gitcmd.GitCommands,
		ghubcmd.CmdForGithub,
		glabcmd.GitLabCmd,
		httpcmd.HttpCmd,
		gocmd.GoToolsCmd,
		phpcmd.PhpToolsCmd,
		javacmd.JavaToolCmd,
		mdcmd.MkDownCmd,
		pkgmanage.ManageCmd,
		strcmd.StringCmd,
		self.KiteManage,
		taskx.TaskManage,
		sqlcmd.SQLCmd,
		toolcmd.ToolsCmd,
		toolcmd.RunScripts,
		builtin.GenAutoComplete(),
	)

	// app.Add(filewatcher.FileWatcher(nil))
	app.Add(pacutil.PacTools.WithHidden())

	app.AddAliases("app:init", "init")
	app.AddAliases("app:info", "info")
	app.AddAliases("app:config", "conf", "config")
}

func addListener(cli *gcli.App) {
	cli.On(events.OnAppInitAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Info("kite cli app init completed")
		return
	})

	cli.On(events.OnCmdRunAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Info("kite cli app cmd: %s run completed", ctx.Cmd.ID())
		return
	})

	cli.On(gcli.EvtCmdNotFound, func(ctx *gcli.HookCtx) bool {
		app.Log().Infof("kite cli app event: %s, TODO handle", ctx.Name())

		// TODO
		return false
	})
}
