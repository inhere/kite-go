package bootstrap

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/appconst"
	"github.com/inhere/kite/internal/initlog"
)

var defaultBaseDir string

// MustRun boot and run app
func MustRun(ka *app.KiteApp) {
	if err := Boot(ka); err != nil {
		cliutil.Errorp(" ERROR ")
		fmt.Println(err)
		os.Exit(-1)
	}
	// goutil.MustOK(err)

	// to run
	app.Run()
}

// Boot app
func Boot(ka *app.KiteApp) error {
	defaultBaseDir = sysutil.ExpandPath(sysutil.Getenv(appconst.EnvKiteBaseDir, appconst.KiteDefaultBaseDir))

	ka.AddPreLoader(BootEnv, func(ka *app.KiteApp) error {
		return initlog.Init(appconst.EnvInitLogLevel)
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
