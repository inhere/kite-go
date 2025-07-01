package bootstrap

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/appconst"
	"github.com/inhere/kite-go/internal/initlog"
)

var defaultBaseDir string

// MustRun boot and run app
func MustRun(ka *app.KiteApp) {
	// boot app
	MustBoot(ka)

	// to run
	app.Run()
}

// MustBoot boot app, if it has error will exit
func MustBoot(ka *app.KiteApp) {
	if err := Boot(ka); err != nil {
		cliutil.Errorp(" ERROR ")
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Boot app
func Boot(ka *app.KiteApp) error {
	defaultBaseDir = sysutil.ExpandPath(sysutil.Getenv(appconst.EnvKiteBaseDir, appconst.KiteDefaultBaseDir))

	ka.AddPreLoader(BootEnv, func(ka *app.KiteApp) error {
		return initlog.Init()
	})

	ka.AddBootFuncs(
		BootConfig,
		BootAppInfo,
		BootLogger,
		BootSrvLogger,
		BootI18n,
		BootCli,
	)
	addServiceBoot(ka)

	return ka.Boot()
}
