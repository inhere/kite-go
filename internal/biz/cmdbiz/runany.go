package cmdbiz

import (
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cliutil/cmdline"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/sysutil/cmdr"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/initlog"
	"github.com/inhere/kite-go/pkg/kscript"
)

// Kas kite command alias map data
var Kas maputil.Aliases

// RunAny handle.
// will try alias, script, plugin, system-cmd ...
func RunAny(name string, args []string, ctx *kscript.RunCtx) error {
	// maybe is kite command alias
	if Kas.HasAlias(name) {
		initlog.L.Infof("TIP: %q is an cli command alias, will run it with %v\n", name, args)
		return RunKiteCmdByAlias(name, args)
	}

	ctx = kscript.EnsureCtx(ctx)
	ctx.BeforeFn = func(si *kscript.ScriptInfo, ctx *kscript.RunCtx) {
		initlog.L.Infof("TIP: %q is a script name, will run it with %v\n", name, args)
		show.AList("Script Info", si)
		show.AList("Run Context", ctx)
	}

	if !ctx.Verbose {
		ctx.BeforeFn = nil
	}

	// try run as script/script-file
	found, err := app.Scripts.TryRun(name, args, ctx)
	if found {
		return err
	}

	// TODO is plugin

	// maybe is system command name
	if sysutil.HasExecutable(name) {
		initlog.L.Infof("TIP: %q is a executable file on system, will run it with %v\n", name, args)
		return cmdr.NewCmd(name, args...).FlushRun()
	}
	return errorx.Rawf("%q is not an alias OR script OR plugin OR system command name", name)
}

// RunKiteCmdByAlias handle
func RunKiteCmdByAlias(name string, inArgs []string) error {
	if !Kas.HasAlias(name) {
		return errorx.Newf("kite alias command %q is not found", name)
	}

	str := Kas.ResolveAlias(name)
	clp := cmdline.NewParser(str)

	cmd, args := clp.BinAndArgs()
	if len(inArgs) > 0 {
		args = append(args, inArgs...)
	}

	if !app.Cli().HasCommand(cmd) {
		return errorx.Rawf("cli command %q not exist, but config in 'aliases.%s'", cmd, name)
	}
	return app.Cli().RunCmd(cmd, args)
}
