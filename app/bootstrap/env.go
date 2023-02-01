package bootstrap

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/ini/v2/dotenv"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

// BootEnv config for kite
func BootEnv(ka *app.KiteApp) error {
	if dotenvFile := findDotEnvFile(ka); dotenvFile != "" {
		initlog.L.Info("find and load kite .env file:", dotenvFile)
		ka.SetDotenvFile(dotenvFile)

		if err := dotenv.LoadFiles(dotenvFile); err != nil {
			return err
		}
	}

	return nil
}

// findDotEnvFile find .env config file
func findDotEnvFile(ka *app.KiteApp) string {
	fileName := appconst.DotEnvFileName
	confFile := envutil.Getenv(appconst.EnvKiteDotEnv, defaultBaseDir+"/"+fileName)

	maybeFiles := []string{
		confFile,
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
