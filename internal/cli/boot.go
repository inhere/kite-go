package cli

import (
	"os"
	"time"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/appconst"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/internal/cli/aicmd"
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
		aicmd.AICommand,
		devcmd.DevToolsCmd,
		fscmd.FsCmd,
		gitcmd.GitCommands,
		ghubcmd.GithubCmd,
		glabcmd.GitLabCmd,
		httpcmd.HttpCmd,
		syscmd.SysCmd,
		appcmd.ManageCmd,
		// extcmd.UserExtCmd,
		textcmd.TextToolCmd,
		jsoncmd.JSONToolCmd,
		// extcmd.XFileCmd,
		toolcmd.ToolsCmd,
		toolcmd.RunAnyCmd,
		toolcmd.NewKScriptCmd(),
		extcmd.PlugCmd,
		builtin.GenAutoComplete().WithHidden(),
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

var (
	autoDir string
	workdir string
	waitSec int
	// set workdir by env KITE_WORKDIR
	defWorkdir = sysutil.Getenv(appconst.EnvKiteWorkdir)
)

func addListener(cli *gcli.App) {
	cli.On(events.OnAppInitAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().WithValue("workdir", cli.WorkDir()).Info("kite cli app init completed. osArgs:", os.Args[1:])
		if err := changeWorkdir(cli, defWorkdir); err != nil {
			colorp.Redln(err.Error())
		}
		return
	})

	// bind new app options
	cli.On(events.OnAppBindOptsAfter, onAppBindOptsAfter(cli))

	cli.On(events.OnCmdRunBefore, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().
			WithValue("workdir", cli.WorkDir()).
			Infof("%s: will run the command %q with args: %v", ctx.Name(), ctx.Cmd.ID(), ctx.Cmd.RawArgs())
		cmdbiz.ProxyCC.AutoSetByCmd(ctx.Cmd)
		return
	})

	cli.On(events.OnCmdRunAfter, func(ctx *gcli.HookCtx) (stop bool) {
		app.Log().Infof("%s: kite cli app command %q run completed", ctx.Name(), ctx.Cmd.ID())
		return
	})

	cli.On(events.OnCmdNotFound, onCmdNotFound)

	cli.On(events.OnAppExit, func(ctx *gcli.HookCtx) (stop bool) {
		if waitSec > 0 {
			app.Log().Infof("%s: will wait %d seconds before app exit. code=%d", ctx.Name(), waitSec, ctx.Int("code"))
			time.Sleep(time.Duration(waitSec) * time.Second)
		}
		return
	})

}

func onAppBindOptsAfter(cli *gcli.App) gcli.HookFunc {
	return func(ctx *gcli.HookCtx) (stop bool) {
		cli.Flags().IntOpt2(&waitSec, "wait", "wait some `seconds` after run command")

		cli.Flags().StrOpt2(&autoDir, "auto-dir,auto-chdir", "auto find dir by name and change workdir",
			gflag.WithValidator(func(val string) error {
				if val == "" {
					return nil
				}

				relDir, changed := fsutil.SearchNameUpx(cli.WorkDir(), val)
				if changed {
					goutil.MustOK(cli.ChWorkDir(relDir))
					cliutil.Yellowf("NOTICE: auto founded dirname %q and will chdir to: %s\n", val, relDir)
				}
				return nil
			}))

		cli.Flags().StrOpt2(&workdir, "workdir,w", "set workdir for run app command",
			gflag.WithValidator(func(val string) error {
				return changeWorkdir(cli, val)
			}),
		)
		return false
	}
}

func onCmdNotFound(ctx *gcli.HookCtx) (stop bool) {
	name := ctx.Str("name")
	args := ctx.Strings("args")
	app.Log().
		WithValue("workdir", ctx.App.WorkDir()).
		WithValue("args", args).
		Infof("%s: handle kite cli command not found: %s", ctx.Name(), name)

	// ctx := &kscript.RunCtx{}
	if err := cmdbiz.RunAny(name, args, nil); err != nil {
		colorp.Warnln("RunAny Error >", err)
	}
	stop = true
	return
}

func changeWorkdir(cli *gcli.App, val string) error {
	if val == "" {
		return nil
	}
	if !fsutil.DirExist(val) {
		return errorx.Err("The workdir not exists: " + val)
	}

	goutil.MustOK(cli.ChWorkDir(val))
	colorp.Yellowf("NOTICE: set app workdir to: %s\n", val)
	return nil
}
