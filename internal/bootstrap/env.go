package bootstrap

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/ini/v2/dotenv"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/appconst"
)

// BootEnv config for kite
func BootEnv(ka *app.KiteApp) error {
	if dotenvFile := findDotEnvFile(ka); dotenvFile != "" {
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
		ka.BinDir() + "/" + fileName,
		ka.WorkDir() + "/" + fileName,
	}

	for _, file := range maybeFiles {
		if fsutil.IsFile(file) {
			return file
		}
	}
	return ""
}
