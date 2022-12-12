package bootstrap

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/i18n"
	"github.com/gookit/slog"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/internal/initlog"
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

func BootI18n(ka *app.KiteApp) error {
	initlog.L.Info("load and init language config files")

	langConf := app.Cfg().SubDataMap("language")

	// lang
	langDir := "resource/language"
	if fsutil.IsDir(langDir) {
		i18n.Init(langDir, langConf.Str("defLang"), langConf.StringMap("langMap"))
	}

	return nil
}

func BootApp(ka *app.KiteApp) error {
	initlog.L.Info("init kite application info config")

	ka.Info = &app.Info{
		Branch:    kite.Branch,
		Version:   kite.Version,
		Revision:  kite.Revision,
		GoVersion: kite.GoVersion,
		PublishAt: kite.PublishAt,
		UpdatedAt: kite.UpdatedAt,
	}

	return nil
}
