package bootstrap

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/i18n"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/initlog"
)

// BootI18n info
func BootI18n(ka *app.KiteApp) error {

	langConf := app.Cfg().SubDataMap("language")

	// lang
	langDir := ka.PathResolve(langConf.Str("lang_dir"))
	if fsutil.IsDir(langDir) {
		initlog.L.Infof("load and init language config files in dir: %s", langDir)
		i18n.Init(langDir, langConf.Str("def_lang"), langConf.StringMap("lang_map"))
	}

	return nil
}

// BootAppInfo info
func BootAppInfo(ka *app.KiteApp) error {
	ka.Info = &app.Info{
		Branch:    kite.Branch,
		Version:   kite.Version,
		Revision:  kite.Revision,
		GoVersion: kite.GoVersion,
		PublishAt: kite.PublishAt,
		UpdatedAt: kite.UpdatedAt,
	}

	initlog.L.Info("init kite application info config:", structs.ToString(ka.Info))

	return nil
}
