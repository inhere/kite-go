package app

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/internal/appconst"
	"github.com/inherelab/kite/internal/initlog"
)

// Init config
func Init(app *KiteApp) error {
	initlog.L.Info("init kite application, info, main config")
	app.Info = &Info{
		Branch:    kite.Branch,
		Version:   kite.Version,
		Revision:  kite.Revision,
		GoVersion: kite.GoVersion,
		PublishAt: kite.PublishAt,
		UpdatedAt: kite.UpdatedAt,
	}

	confFile := findConfFile(app)
	if confFile == "" {
		return nil
	}

	initlog.L.Info("load main config file:", confFile)
	app.mainFile = confFile
	err := app.cfg.LoadFiles(confFile)
	if err != nil {
		return err
	}

	// map main config
	err = app.cfg.MapOnExists(appconst.ConfKeyApp, app.Conf)
	if err != nil {
		return err
	}

	return app.init()
}

// findConfFile find main config file
func findConfFile(app *KiteApp) string {
	file := envutil.Getenv(appconst.EnvKiteConfig, sysutil.ExpandPath(appconst.KiteDefaultConfigFile))
	if fsutil.IsFile(file) {
		return file
	}

	file = app.WorkDir() + "/" + appconst.KiteConfigName
	if fsutil.IsFile(file) {
		return file
	}

	file = app.BinDir() + "/" + appconst.KiteConfigName
	if fsutil.IsFile(file) {
		return file
	}
	return ""
}
