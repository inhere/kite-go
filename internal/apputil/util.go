package apputil

import (
	"os"

	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inherelab/kite/internal/appconst"
)

// SetEnvs to os
func SetEnvs(mp map[string]string) {
	for key, value := range mp {
		_ = os.Setenv(key, value)
	}
}

// FindConfFile find main config file
func FindConfFile() string {
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
