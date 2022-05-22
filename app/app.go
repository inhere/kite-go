package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
	"github.com/inherelab/kite"
)

// Env names
const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvPre   = "pre"
	EnvProd  = "prod"
)

// Info for kite app
var Info = &struct {
	Branch    string
	Version   string
	PubDate   string
	Revision  string
	BuildDate string
	GoVersion string
}{
	Version: "1.0.0",
	PubDate: "2021-02-14 13:14",
}

var (
	KiteMode = envutil.Getenv("KITE_DEBUG", slog.DebugLevel.LowerName())
	KiteConf = envutil.Getenv("KITE_CONFIG", sysutil.UserDir(".config/kite.yml"))
)

// IsDebug mode
func IsDebug() bool {
	return slog.LevelByName(KiteMode) >= slog.DebugLevel
}

// BootLoader for app start boot
type BootLoader interface {
	// Boot do something before application run
	Boot(app *App) error
}

// BootFunc for application
type BootFunc func(app *App) error

// Boot do something
func (fn BootFunc) Boot(app *App) error {
	return fn(app)
}

// App struct
type App struct {
	*gcli.App
	// config for app
	cfg *config.Config

	// main config file path.
	cfgFile string
}

func A() *App {
	return nil
}

// C cfg get the config.Config
func C() *config.Config {
	return nil
}

// Rux get the web router
func Rux() *rux.Router {
	return nil
}

func Run() {

}

func NewApp() *App {
	ca := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"

		a.Version = kite.Version
	})

	return &App{
		App: ca,
	}
}
