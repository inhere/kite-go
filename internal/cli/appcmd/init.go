package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/appconst"
)

var initDotenv = `
# custom env settings
KITE_VERBOSE = debug
KITE_INIT_LOG = debug
KITE_GLAB_HOST = http://gitlab.your.com

`

var ikOpts = struct {
	dryRun bool
}{}

// KiteInitCmd command
var KiteInitCmd = &gcli.Command{
	Name: "init",
	Desc: "init kite .env, config to the user home dir",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&ikOpts.dryRun, "dry-run, dry", "run workflow, dont real execute")
	},
	Func: func(c *gcli.Command, args []string) error {
		uHome := sysutil.UserHomeDir()
		c.Infoln("Found user home dir:", uHome)

		baseDir := app.App().BaseDir
		dotenvFile := fsutil.JoinPaths(baseDir, appconst.DotEnvFileName)
		c.Infoln("Will init kite .env and config and more...")

		if ikOpts.dryRun {
			c.Warnln("TIP: on DRY-RUN mode, will not be execute any operation")
		}

		qr := goutil.NewQuickRun()
		qr.Add(func(ctx *structs.Data) error {
			c.Infoln("- make the base dir:", baseDir)
			if ikOpts.dryRun {
				return nil
			}

			return fsutil.Mkdir(baseDir, fsutil.DefaultDirPerm)
		}, func(ctx *structs.Data) error {
			c.Infoln("- init the .env file:", dotenvFile)
			if ikOpts.dryRun {
				return nil
			}

			_, err := fsutil.PutContents(dotenvFile, initDotenv, fsutil.FsCWTFlags)
			return err
		}, func(ctx *structs.Data) error {
			c.Infof("- init data dirs in %s: scripts, plugins\n", baseDir)
			if ikOpts.dryRun {
				return nil
			}

			return fsutil.MkSubDirs(fsutil.DefaultDirPerm, baseDir, "scripts", "plugins")
		})

		return qr.Run()
	},
}
