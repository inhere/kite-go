package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/internal/appconst"
	"github.com/inherelab/kite/internal/initlog"
)

// Env names
const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvTest  = "test"
	EnvPre   = "pre"
	EnvProd  = "prod"
)

var (
	KiteVerbose = envutil.Getenv(appconst.EnvKiteVerbose, slog.WarnLevel.LowerName())
)

// IsDebug mode
func IsDebug() bool {
	return slog.LevelByName(KiteVerbose) >= slog.DebugLevel
}

// BootLoader for app start boot
type BootLoader interface {
	// Boot do something before application run
	Boot(app *KiteApp) error
}

// BootFunc for application
type BootFunc func(app *KiteApp) error

// Boot do something
func (fn BootFunc) Boot(app *KiteApp) error {
	return fn(app)
}

// Info for kite app
type Info struct {
	Branch    string
	Version   string
	Revision  string
	GoVersion string
	BuildDate string
	PublishAt string
	UpdatedAt string
}

// KiteApp kite app struct
type KiteApp struct {
	*Conf
	*Info
	*gcli.App
	// config data for app
	cfg *config.Config
}

// IsDebug mode
func (app *KiteApp) IsDebug() bool {
	return IsDebug()
}

// Cfg instance
func (app *KiteApp) Cfg() *config.Config {
	return app.cfg
}

func (app *KiteApp) init() error {
	app.ensurePaths()

	return app.LoadIncludeConfigs()
}

// LoadIncludeConfigs from conf.IncludeConfig
func (app *KiteApp) LoadIncludeConfigs() error {
	ln := len(app.IncludeConfig)
	if ln == 0 {
		return nil
	}

	initlog.L.Info("load include config files from 'include_config'")
	filePaths := make([]string, 0, ln)
	for _, file := range app.IncludeConfig {
		if len(file) < 2 {
			continue
		}

		// is relative path
		if file[0] != OSPathSepChar && file[0] != PathAliasPrefix {
			filePaths = append(filePaths, app.ConfigPath(file))
		} else {
			filePaths = append(filePaths, app.PathResolve(file))
		}
	}

	return app.cfg.LoadFiles(filePaths...)
}

var kiteApp = newInitApp()

// App instance
func App() *KiteApp {
	return kiteApp
}

// Run boot and run app
func Run() {
	kiteApp.Run(nil)
}

// Cfg get the config.Config
func Cfg() *config.Config {
	return kiteApp.cfg
}

// Rux get the web router
func Rux() *rux.Router {
	return nil
}

func newInitApp() *KiteApp {
	cliApp := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"

		a.Version = kite.Version
	})

	app := &KiteApp{
		App: cliApp,
		// Info: info,
		Conf: newDefaultConf(),
	}

	app.cfg = config.NewWith("kite", func(c *config.Config) {
		c.AddDriver(yamlv3.Driver)
		c.WithOptions(func(opt *config.Options) {
			opt.DecoderConfig.TagName = "json"
		})
	})

	err := Init(app)
	if err != nil {
		panic(err)
	}

	return app
}
