package bootstrap

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/internal/appconst"
	"github.com/inherelab/kite/internal/initlog"
)

func BootEnv(ka *app.KiteApp) error {
	confFile := findConfFile(ka)
	initlog.L.Info("load main config file:", confFile)
	ka.SetConfFile(confFile)
	return nil
}

// findConfFile find main config file
func findConfFile(ka *app.KiteApp) string {
	file := envutil.Getenv(appconst.EnvKiteConfig, sysutil.ExpandPath(appconst.KiteDefaultConfigFile))
	if fsutil.IsFile(file) {
		return file
	}

	file = ka.WorkDir() + "/" + appconst.KiteConfigName
	if fsutil.IsFile(file) {
		return file
	}

	file = ka.BinDir() + "/" + appconst.KiteConfigName
	if fsutil.IsFile(file) {
		return file
	}
	return ""
}
