package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/initlog"
)

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
