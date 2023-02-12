package cli

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/gookit/gcli/v3/events"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/cli/appcmd"
	"github.com/inhere/kite/internal/cli/devcmd"
	"github.com/inhere/kite/internal/cli/fscmd"
	"github.com/inhere/kite/internal/cli/ghubcmd"
	"github.com/inhere/kite/internal/cli/gitcmd"
	"github.com/inhere/kite/internal/cli/glabcmd"
	"github.com/inhere/kite/internal/cli/httpcmd"
	"github.com/inhere/kite/internal/cli/pkgcmd"
	"github.com/inhere/kite/internal/cli/strcmd"
	"github.com/inhere/kite/internal/cli/syscmd"
	"github.com/inhere/kite/internal/cli/taskx"
	"github.com/inhere/kite/internal/cli/toolcmd"
	"github.com/inhere/kite/pkg/pacutil"
)

// Boot commands to gcli.App
func Boot(cli *gcli.App) {
	addListener(cli)

	addCommands(cli)

	addAliases(cli)
}

// addCommands commands to gcli.App
func addCommands(cli *gcli.App) {
	cli.Add(
		devcmd.DevToolsCmd,
		fscmd.FsCmd,
		gitcmd.GitCommands,
		ghubcmd.GithubCmd,
		glabcmd.GitLabCmd,
		httpcmd.HttpCmd,
		pkgcmd.PkgManageCmd,
		strcmd.StringCmd,
		syscmd.SysCmd,
		appcmd.ManageCmd,
		taskx.TaskManage,
		toolcmd.ToolsCmd,
		toolcmd.RunAnyCmd,
		builtin.GenAutoComplete(func(c *gcli.Command) {
			c.Hidden = true
		}),
	)

	// app.Add(filewatcher.FileWatcher(nil))
	cli.Add(pacutil.PacTools.WithHidden())
}

func addAliases(cli *gcli.App) {
	// built in alias
	cli.AddAliases("app:init", "init")
	cli.AddAliases("app:info", "info")
	cli.AddAliases("app:config", "conf", "config")
}

func addListener(cli *gcli.App) {
	cli.On(events.OnAppInitAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Info("kite cli app init completed")
		return
	})

	cli.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Infof("will run the command: %s with args: %v", ctx.Cmd.ID(), ctx.Cmd.RawArgs())
		return
	})

	cli.On(events.OnCmdRunAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Infof("kite cli app command: %s run completed", ctx.Cmd.ID())
		return
	})

	cli.On(events.OnCmdNotFound, cmdNotFund)
}

func cmdNotFund(ctx *gcli.HookCtx) (stop bool) {
	name := ctx.Str("name")
	app.Log().Infof("kite cli event: %s, not found: %s", ctx.Name(), name)

	color.Infoln("fire not found event:", ctx.Name(),
		", name:", name, ",args:", ctx.Strings("args"),
	)

	// TODO runAny.Run()
	return
}
