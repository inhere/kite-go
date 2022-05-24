package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/internal/appconst"
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
	KiteRMode = envutil.Getenv(appconst.EnvKiteVerbose, slog.WarnLevel.LowerName())
)

// IsDebug mode
func IsDebug() bool {
	return slog.LevelByName(KiteRMode) >= slog.DebugLevel
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

// Config struct
//
// Gen by:
//   kite go gen st -s @c -t json --name AppConfig
type Config struct {
	// BaseDir base dir
	BaseDir string `json:"base_dir"`
	// TmpDir tmp dir
	TmpDir string `json:"tmp_dir"`
	// CacheDir cache dir
	CacheDir string `json:"cache_dir"`
	// ConfigDir config dir
	ConfigDir string `json:"config_dir"`
	// ResourceDir resource dir
	ResourceDir string `json:"resource_dir"`
	// IncludeConfig include config files
	IncludeConfig []string `json:"include_config"`
}

// KiteApp kite app struct
type KiteApp struct {
	*Info
	*Config
	*gcli.App

	// the main config file path.
	cfgFile string
	// config for app
	cfg *config.Config

	WorkDir string `json:"work_dir"`
}

// CfgFile get main config file
func (app *KiteApp) CfgFile() string {
	return app.cfgFile
}

// IsDebug mode
func (app *KiteApp) IsDebug() bool {
	return IsDebug()
}

// Cfg instance
func (app *KiteApp) Cfg() *config.Config {
	return app.cfg
}

var kiteApp = newInitApp()

var Aliases = &structs.Aliases{}

// App instance
func App() *KiteApp {
	return kiteApp
}

// Run boot and run app
func Run() {
	kiteApp.Run(nil)
}

// C get the config.Config
func C() *config.Config {
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

	cfg := config.NewWith("kite", func(c *config.Config) {
		c.AddDriver(yamlv3.Driver)
		c.WithOptions(func(opt *config.Options) {
			opt.DecoderConfig.TagName = "json"
		})
	})

	app := &KiteApp{
		cfg: cfg,
		App: cliApp,
		// Info: info,
		Config: &Config{},
	}

	err := Init(app)
	if err != nil {
		panic(err)
	}

	return app
}
