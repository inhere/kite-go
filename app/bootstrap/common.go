package bootstrap

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/i18n"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/initlog"
)

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
