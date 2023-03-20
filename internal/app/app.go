package app

import (
	"os"
	"sync"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/internal/appconst"
	"github.com/inhere/kite-go/internal/initlog"
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
	KiteVerbose = sysutil.Getenv(appconst.EnvKiteVerbose, slog.WarnLevel.LowerName())
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

	// pre-bootloaders
	preLoaders []BootLoader
	// app bootloaders
	bootloaders []BootLoader
	// shutdown callback functions
	shutdown []func()
}

// InitPaths for app
func (ka *KiteApp) InitPaths() {
	ka.ensurePaths()
}

// AddPreLoader to app
func (ka *KiteApp) AddPreLoader(bfs ...BootFunc) *KiteApp {
	for _, bf := range bfs {
		ka.preLoaders = append(ka.preLoaders, bf)
	}
	return ka
}

// AddBootFuncs to app
func (ka *KiteApp) AddBootFuncs(bfs ...BootFunc) *KiteApp {
	for _, bf := range bfs {
		ka.bootloaders = append(ka.bootloaders, bf)
	}
	return ka
}

// AddLoaders to app
func (ka *KiteApp) AddLoaders(bls ...BootLoader) *KiteApp {
	ka.bootloaders = append(ka.bootloaders, bls...)
	return ka
}

// AddLoader to app
func (ka *KiteApp) AddLoader(bl BootLoader) *KiteApp {
	ka.bootloaders = append(ka.bootloaders, bl)
	return ka
}

// Boot app start
func (ka *KiteApp) Boot() error {
	err := ka.runBootloaders(ka.preLoaders)
	if err != nil {
		return err
	}

	return ka.runBootloaders(ka.bootloaders)
}

// SetConfFile path.
func (ka *KiteApp) runBootloaders(loaders []BootLoader) error {
	for _, loader := range loaders {
		if bc, ok := loader.(BootChecker); ok {
			if !bc.BeforeBoot() {
				initlog.L.Debugf("skip boot on %v.BeforeBoot() return false", bc)
				continue
			}
		}

		if err := loader.Boot(ka); err != nil {
			return errorx.Wrapf(err, "bootloader run fail on %#v", loader)
		}
	}
	return nil
}

// SetConfFile path.
func (ka *KiteApp) SetConfFile(file string) {
	ka.confFile = file
}

// OnShutdown handler.
func (ka *KiteApp) OnShutdown(fn func()) {
	ka.shutdown = append(ka.shutdown, fn)
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
			Config: &Config{},
		}
	})

	return kiteApp
}

// Run app
func Run() {
	code := Cli().Run(nil)

	for _, fn := range kiteApp.shutdown {
		fn()
	}
	os.Exit(code)
}
