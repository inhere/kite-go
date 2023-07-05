package cli_test

import (
	"testing"

	"github.com/gookit/color/colorp"
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
	colorp.Successln("the kite test application bootstrap success, workdir:", sysutil.Workdir())
	m.Run()
}

func TestApp_run(t *testing.T) {
	st := app.Cli().RunLine("-h")
	assert.Eq(t, st, 0)
}

func TestApp_chdir(t *testing.T) {
	st := app.Cli().RunLine("--chdir .git app info")
	assert.Eq(t, st, 0)
}

func TestCmd_http_tpl_send(t *testing.T) {
	st := app.Cli().RunLine("http tpl-send -d github --api releases-latest -P -t 2000 -v owner=gookit -v repo=slog")
	assert.Eq(t, st, 0)
}
