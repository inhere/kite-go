package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

var defaultBaseDir = sysutil.ExpandPath(appconst.KiteDefaultDataDir)

// MustBoot app
func MustBoot(ka *app.KiteApp) {
	goutil.MustOK(Boot(ka))
}

// Boot app
func Boot(ka *app.KiteApp) error {
	initlog.L.Info("bootstrap the kite application, register boot loaders and run")

	ka.AddBootFuncs(
		BootAppInfo,
		BootEnv,
		BootConfig,
		BootLogger,
		BootI18n,
		BootCli,
	)

	addServiceBoot(ka)

	return ka.Boot()
}
