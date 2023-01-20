package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/gookit/slog"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/initlog"
)

// MustBoot app
func MustBoot(ka *app.KiteApp) {
	goutil.MustOK(Boot(ka))
}

// Boot app
func Boot(ka *app.KiteApp) error {
	initlog.L.Info("bootstrap the kite application, register boot loaders")

	ka.AddBootFuncs(
		BootEnv,
		BootAppInfo,
		BootConfig,
		BootLogger,
		BootI18n,
		BootCli,
	)

	addServiceBoot(ka)

	return ka.Boot()
}

func configSlog() {
	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.WarnLevel

		f := logger.Formatter.(*slog.TextFormatter)
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}}\n")
	})

	slog.SetLogLevel(slog.LevelByName(app.KiteVerbose))
}
