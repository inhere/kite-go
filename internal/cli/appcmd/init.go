package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/appconst"
)

var ikOpts = struct {
	dryRun bool
	force  bool
}{}

// KiteInitCmd command
var KiteInitCmd = &gcli.Command{
	Name: "init",
	Desc: "init kite .env, config to the user home dir",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&ikOpts.dryRun, "dry-run, dry", "run workflow, dont real execute")
		c.BoolOpt2(&ikOpts.force, "force, f", "force re-init kite app config")
	},
	Func: func(c *gcli.Command, args []string) error {
		uHome := sysutil.UserHomeDir()
		c.Infoln("Found user home dir:", uHome)

		baseDir := app.App().BaseDir
		dotenvFile := fsutil.JoinPaths(baseDir, appconst.DotEnvFileName)
		confFile := fsutil.JoinPaths(baseDir, appconst.KiteConfigName)
		c.Infoln("Will init kite .env and config and more...")

		if ikOpts.dryRun {
			c.Warnln("TIP: on DRY-RUN mode, will not be execute any operation")
		}

		qr := goutil.NewQuickRun()
		qr.Add(func(ctx *structs.Data) error {
			c.Infoln("- Make the kite base dir:", baseDir)
			if ikOpts.dryRun {
				return nil
			}

			return fsutil.Mkdir(baseDir, fsutil.DefaultDirPerm)
		}, func(ctx *structs.Data) error {
			c.Infoln("- Init the .env file:", dotenvFile)
			if ikOpts.dryRun {
				return nil
			}

			if !ikOpts.force && fsutil.IsFile(dotenvFile) {
				c.Warnln("  Exists, skip write!")
				return nil
			}

			text := byteutil.SafeString(kite.EmbedFs.ReadFile(".env.example"))
			_, err := fsutil.PutContents(dotenvFile, text, fsutil.FsCWTFlags)
			return err
		}, func(ctx *structs.Data) error {
			c.Infoln("- Init the main config file:", confFile)
			if ikOpts.dryRun {
				return nil
			}

			if !ikOpts.force && fsutil.IsFile(confFile) {
				c.Warnln("  Exists, skip write!")
				return nil
			}

			text := byteutil.SafeString(kite.EmbedFs.ReadFile("kite.example.yml"))
			_, err := fsutil.PutContents(confFile, text, fsutil.FsCWTFlags)
			return err
		}, func(ctx *structs.Data) error {
			subDirs := []string{"config", "data", "scripts", "plugins"}
			c.Infof("- Init data dirs in %s: %v\n", baseDir, subDirs)
			if ikOpts.dryRun {
				return nil
			}

			return fsutil.MkSubDirs(fsutil.DefaultDirPerm, baseDir, subDirs...)
		}, func(ctx *structs.Data) error {
			c.Infof("- Init kite config files to %s/config\n", baseDir)
			if ikOpts.dryRun {
				return nil
			}

			entries, err := kite.EmbedFs.ReadDir("config")
			if err != nil {
				return err
			}

			for _, entry := range entries {
				path := "config/" + entry.Name()
				fmt.Println("  Init write the", path)

				dstFile := baseDir + "/" + path
				if !ikOpts.force && fsutil.IsFile(dstFile) {
					c.Warnln("   Exists, skip write!")
					continue
				}

				text := byteutil.SafeString(kite.EmbedFs.ReadFile(path))
				_, err = fsutil.PutContents(dstFile, text, fsutil.FsCWTFlags)
				if err != nil {
					return err
				}
			}

			return err
		})

		err := qr.Run()
		if err == nil {
			cliutil.Successln(" ✅ Init Completed")
		}
		return err
	},
}
