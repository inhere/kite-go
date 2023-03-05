package bootstrap

import (
	"os"
	"runtime"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

// BootConfig for kite
func BootConfig(ka *app.KiteApp) error {
	cfg := config.NewWith("kite", func(c *config.Config) {
		c.AddDriver(yamlv3.Driver)
		c.WithOptions(config.WithTagName("json"), func(opt *config.Options) {
			opt.ParseEnv = true
			opt.ParseDefault = true
		})
	})

	confFile := findMainConfFile(ka, appconst.KiteConfigName)
	if confFile != "" {
		ka.SetConfFile(confFile)
		initlog.L.Info("load the kite main config file:", confFile)
		if err := cfg.LoadFiles(confFile); err != nil {
			return err
		}
	} else {
		initlog.L.Warn("the kite main config file not found")
	}

	// map app config
	if err := mapAppConfig(ka, cfg); err != nil {
		return err
	}

	// load include configs
	if err := loadIncludeConfigs(ka, cfg); err != nil {
		return err
	}

	if app.IsDebug() {
		dump.P(cfg.Data())
	}

	app.Add(app.ObjConf, cfg)
	return nil
}

// findConfFile find main config file
func findMainConfFile(ka *app.KiteApp, fileName string) string {
	confFile := envutil.Getenv(appconst.EnvKiteConfig, defaultBaseDir+"/"+fileName)
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

// LoadIncludeConfigs from conf.IncludeConfig
func mapAppConfig(ka *app.KiteApp, cfg *config.Config) error {
	err := cfg.MapOnExists(appconst.ConfKeyApp, ka.Config)
	if err != nil {
		return err
	}

	ka.InitPaths()
	initlog.L.Debug("app.Config init ok, kite base_dir is", ka.BaseDir)
	return nil
}

// can use var in filepath
var pathVars = map[string]string{
	"$os":   runtime.GOOS,
	"$user": os.Getenv("USER"),
}

// LoadIncludeConfigs from conf.IncludeConfig
func loadIncludeConfigs(ka *app.KiteApp, cfg *config.Config) error {
	ln := len(ka.IncludeConfig)
	if ln == 0 {
		return nil
	}

	initlog.L.Info("load include config files from 'include_config'", ka.IncludeConfig)
	filePaths := make([]string, 0, ln)
	for _, file := range ka.IncludeConfig {
		if len(file) < 2 {
			continue
		}

		var filePath string
		// contains var
		if strutil.ContainsByte(file, '$') {
			file = strutil.Replaces(file, pathVars)
		}

		// is relative path
		if file[0] != app.OSPathSepChar && !app.IsAliasPath(file) {
			filePath = ka.ConfigPath(file)
		} else {
			filePath = ka.PathResolve(file)
		}

		initlog.L.Debugf("load the include file: %s", filePath)
		filePaths = append(filePaths, filePath)
	}

	// platform config file: config.darwin.yml
	platFile := ka.ConfigPath("config." + runtime.GOOS + ".yml")
	if fsutil.IsFile(platFile) {
		initlog.L.Info("will auto load the platform config file:", platFile)
		filePaths = append(filePaths, platFile)
	}

	return cfg.LoadFiles(filePaths...)
}
