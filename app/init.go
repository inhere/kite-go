package app

import (
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/internal/appconst"
	"github.com/inherelab/kite/internal/initlog"
)

// Init config
func Init(app *KiteApp) error {
	initlog.L.Info("init kite application")
	info := &Info{
		Branch:    kite.Branch,
		Version:   kite.Version,
		Revision:  kite.Revision,
		GoVersion: kite.GoVersion,
		PublishAt: kite.PublishAt,
		UpdatedAt: kite.UpdatedAt,
	}

	app.Info = info

	confFile := findConfFile()
	if confFile == "" {
		return nil
	}

	initlog.L.Info("load main config file:", confFile)
	app.cfgFile = confFile
	err := app.cfg.LoadFiles(confFile)
	if err != nil {
		return err
	}

	// map config
	err = app.cfg.MapOnExists(appconst.ConfKeyApp, app.Config)
	return err
}

// findConfFile find main config file
func findConfFile() string {
	confFile := envutil.Getenv(appconst.EnvKiteConfig, sysutil.UserDir(".kite/"+appconst.KiteConfigFile))
	if fsutil.IsFile(confFile) {
		return confFile
	}

	confFile = cliutil.Workdir() + "/" + appconst.KiteConfigFile
	if fsutil.IsFile(confFile) {
		return confFile
	}

	confFile = cliutil.BinDir() + "/" + appconst.KiteConfigFile
	if fsutil.IsFile(confFile) {
		return confFile
	}
	return ""
}
