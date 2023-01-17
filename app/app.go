package app

import (
	"sync"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/slog"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
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

// KiteApp kite app container
type KiteApp struct {
	*Info
	*Config
	*gcli.Context

	loaders []BootLoader
}

// InitPaths for app
func (ka *KiteApp) InitPaths() {
	ka.ensurePaths()
}

// AddBootFuncs to app
func (ka *KiteApp) AddBootFuncs(bfs ...BootFunc) {
	for _, bf := range bfs {
		ka.loaders = append(ka.loaders, bf)
	}
}

// AddLoaders to app
func (ka *KiteApp) AddLoaders(bls ...BootLoader) {
	ka.loaders = append(ka.loaders, bls...)
}

// Boot app
func (ka *KiteApp) Boot() error {
	for _, loader := range ka.loaders {
		if bc, ok := loader.(BootChecker); ok {
			if !bc.BeforeBoot() {
				initlog.L.Debugf("skip boot on %v.BeforeBoot() return false", bc)
				continue
			}
		}

		if err := loader.Boot(ka); err != nil {
			return errorx.Wrap(err, "boot loader fail on "+goutil.FuncName(loader))
		}
	}
	return nil
}

// Run app
func (ka *KiteApp) Run() {
	Cli().Run(nil)
}

func (ka *KiteApp) SetConfFile(file string) {
	ka.confFile = file
}

var initKa sync.Once
var kiteApp *KiteApp

// App instance
func App() *KiteApp {
	// init app at once
	initKa.Do(func() {
		kiteApp = &KiteApp{
			Context: gcli.GCtx(),
			// Info: info,
			Config: newDefaultConf(),
		}
	})

	return kiteApp
}

// Run boot and run app
func Run() {
	kiteApp.Run()
}
