package bootstrap

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/dump"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

func BootConfig(ka *app.KiteApp) error {
	cfg := config.NewWith("kite", func(c *config.Config) {
		c.AddDriver(yamlv3.Driver)
		c.WithOptions(func(opt *config.Options) {
			opt.ParseEnv = true
			opt.DecoderConfig.TagName = "json"
		})
	})

	confFile := ka.ConfFile()
	if confFile != "" {
		initlog.L.Info("load and init kite config file:", confFile)
		if err := cfg.LoadFiles(confFile); err != nil {
			return err
		}
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

	app.Add("config", cfg)
	return nil
}

// LoadIncludeConfigs from conf.IncludeConfig
func mapAppConfig(ka *app.KiteApp, cfg *config.Config) error {
	err := cfg.MapOnExists(appconst.ConfKeyApp, ka.Config)
	if err != nil {
		return err
	}

	ka.InitPaths()
	return nil
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

		// is relative path
		if file[0] != app.OSPathSepChar && file[0] != app.PathAliasPrefix {
			filePaths = append(filePaths, ka.ConfigPath(file))
		} else {
			filePaths = append(filePaths, ka.PathResolve(file))
		}
	}

	return cfg.LoadFiles(filePaths...)
}
