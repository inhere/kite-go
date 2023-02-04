package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

var defaultBaseDir = sysutil.ExpandPath(appconst.KiteDefaultDataDir)

// MustRun boot and run app
func MustRun(ka *app.KiteApp) {
	goutil.MustOK(Boot(ka))

	// to run
	app.Run()
}

// Boot app
func Boot(ka *app.KiteApp) error {
	ka.AddPreLoader(BootEnv, func(ka *app.KiteApp) error {
		return initlog.Init(appconst.EnvInitLogLevel)
	})

	ka.AddBootFuncs(
		BootAppInfo,
		BootConfig,
		BootLogger,
		BootI18n,
		BootCli,
	)
	addServiceBoot(ka)

	return ka.Boot()
}
