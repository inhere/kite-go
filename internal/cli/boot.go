package cli

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/cliutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/internal/cli/appcmd"
	"github.com/inhere/kite-go/internal/cli/devcmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/jsoncmd"
	"github.com/inhere/kite-go/internal/cli/extcmd"
	"github.com/inhere/kite-go/internal/cli/fscmd"
	"github.com/inhere/kite-go/internal/cli/gitcmd"
	"github.com/inhere/kite-go/internal/cli/gitcmd/ghubcmd"
	"github.com/inhere/kite-go/internal/cli/gitcmd/glabcmd"
	"github.com/inhere/kite-go/internal/cli/httpcmd"
	"github.com/inhere/kite-go/internal/cli/syscmd"
	"github.com/inhere/kite-go/internal/cli/taskcmd"
	"github.com/inhere/kite-go/internal/cli/textcmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd"
	"github.com/inhere/kite-go/pkg/pacutil"
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
		syscmd.SysCmd,
		appcmd.ManageCmd,
		extcmd.UserExtCmd,
		taskcmd.TaskManageCmd,
		textcmd.TextOperateCmd,
		jsoncmd.JSONToolCmd,
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
		app.Log().Infof("%s: will run the command %q with args: %v", ctx.Name(), ctx.Cmd.ID(), ctx.Cmd.RawArgs())
		cmdbiz.ProxyCC.AutoSetByCmd(ctx.Cmd)
		return
	})

	cli.On(events.OnCmdRunAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Infof("%s: kite cli app command %q run completed", ctx.Name(), ctx.Cmd.ID())
		return
	})

	cli.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) (stop bool) {
		name := ctx.Str("name")
		args := ctx.Strings("args")
		app.Log().Infof("%s: handle kite cli command not found: %s", ctx.Name(), name)

		if err := cmdbiz.RunAny(name, args, nil); err != nil {
			cliutil.Warnln("RunAny Error >", err)
		}
		stop = true
		return
	})
}
