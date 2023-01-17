package bootstrap

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/ini/v2/dotenv"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

func BootEnv(ka *app.KiteApp) error {
	if envFile := findConfFile(ka, ".env"); envFile != "" {
		initlog.L.Info("find and load ENV config file:", envFile)

		if err := dotenv.LoadFiles(envFile); err != nil {
			return err
		}
	}

	if confFile := findConfFile(ka, appconst.KiteConfigName); confFile != "" {
		initlog.L.Info("find main config file:", confFile)
		ka.SetConfFile(confFile)
	}

	return nil
}

// findConfFile find main config file
func findConfFile(ka *app.KiteApp, fileName string) string {
	envFile := envutil.Getenv(appconst.EnvKiteConfig, sysutil.ExpandPath(appconst.KiteDefaultDataDir)+"/"+fileName)
	maybeFiles := []string{
		envFile,
		ka.WorkDir() + "/" + fileName,
		ka.BinDir() + "/" + fileName,
	}

	for _, file := range maybeFiles {
		if fsutil.IsFile(file) {
			return file
		}
	}
	return ""
}
