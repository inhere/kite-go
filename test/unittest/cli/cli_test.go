package cli_test

import (
	"testing"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/bootstrap"
	"github.com/inhere/kite-go/internal/initlog"
)

func TestMain(m *testing.M) {
	app.App().AfterPreFn = func(ka *app.KiteApp) error {
		initlog.SetLevel(slog.DebugLevel)
		return nil
	}

	bootstrap.MustBoot(app.App())

	// set verbose level
	gcli.GOpts().Verbose = gcli.VerbDebug
	colorp.Successln("the kite test application bootstrap success, workdir:", sysutil.Workdir())
	m.Run()
}

func TestApp_run(t *testing.T) {
	st := app.Cli().RunLine("-h")
	assert.Eq(t, st, 0)
}

func TestApp_chdir(t *testing.T) {
	st := app.Cli().RunLine("--auto-dir .git app info")
	assert.Eq(t, st, 0)
}

func TestApp_chdir_gitcmd(t *testing.T) {
	st := app.Cli().RunLine("git --auto-root status")
	assert.Eq(t, st, 0)
}
