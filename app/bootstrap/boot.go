package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/gookit/slog"
	"github.com/inhere/kite/app"
)

// MustBoot app
func MustBoot(ka *app.KiteApp) {
	goutil.MustOK(Boot(ka))
}

// Boot app
func Boot(ka *app.KiteApp) error {
	slog.SetLogLevel(slog.LevelByName(app.KiteVerbose))
	slog.Info("bootstrap the kite application, register boot loaders")

	ka.AddLoaders(
		app.BootFunc(BootLogger),
	)

	ka.AddBootFuncs(BootEnv, BootApp, BootConfig, BootI18n, BootCli)

	return ka.Boot()
}
