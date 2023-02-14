package appcmd

import (
	"embed"
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
		c.Warnln("Pre-check:")
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
			idx := ctx.IntVal("index") + 1
			c.Infof("%d. Make the kite base dir: %s\n", idx, baseDir)
			if ikOpts.dryRun {
				return nil
			}

			if !ikOpts.force && fsutil.IsFile(baseDir) {
				c.Warnln("   Exists, skip write!")
				return nil
			}

			return fsutil.Mkdir(baseDir, fsutil.DefaultDirPerm)
		}, func(ctx *structs.Data) error {
			idx := ctx.IntVal("index") + 1
			c.Infof("%d. Init the .env file: %s\n", idx, dotenvFile)
			if ikOpts.dryRun {
				return nil
			}

			if !ikOpts.force && fsutil.IsFile(dotenvFile) {
				c.Warnln("   Exists, skip write!")
				return nil
			}

			text := byteutil.SafeString(kite.EmbedFs.ReadFile(".example.env"))
			_, err := fsutil.PutContents(dotenvFile, text, fsutil.FsCWTFlags)
			return err
		}, func(ctx *structs.Data) error {
			idx := ctx.IntVal("index") + 1
			c.Infof("%d. Init the main config file: %s\n", idx, confFile)
			if ikOpts.dryRun {
				return nil
			}

			if !ikOpts.force && fsutil.IsFile(confFile) {
				c.Warnln("   Exists, skip write!")
				return nil
			}

			text := byteutil.SafeString(kite.EmbedFs.ReadFile("kite.example.yml"))
			_, err := fsutil.PutContents(confFile, text, fsutil.FsCWTFlags)
			return err
		}, func(ctx *structs.Data) error {
			subDirs := []string{"config", "data", "scripts", "plugins", "tmp"}
			idx := ctx.IntVal("index") + 1
			c.Infof("%d. Init subdir in %s: %v\n", idx, baseDir, subDirs)
			if ikOpts.dryRun {
				return nil
			}

			return fsutil.MkSubDirs(fsutil.DefaultDirPerm, baseDir, subDirs...)
		}, func(ctx *structs.Data) error {
			idx := ctx.IntVal("index") + 1
			c.Infof("%d. Init kite config files to %s/config\n", idx, baseDir)
			if ikOpts.dryRun {
				return nil
			}

			return exportEmbedDir(kite.EmbedFs, "config", baseDir+"/config", true)
		})

		cliutil.Warnln("\nStarting init:")
		err := qr.Run()
		if err == nil {
			cliutil.Successln("✅  Init Completed")
		}
		return err
	},
}

func exportEmbedDir(efs embed.FS, dirPath, dstDir string, exportSub bool) error {
	entries, err := efs.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		path := dirPath + "/" + name

		if entry.IsDir() {
			if exportSub {
				err = exportEmbedDir(efs, path, dstDir+"/"+name, exportSub)
			}
			continue
		}

		fmt.Print("   Read and init the ", path)
		dstFile := dstDir + "/" + name
		if !ikOpts.force && fsutil.IsFile(dstFile) {
			cliutil.Warnln("   Exists")
			continue
		}

		text := byteutil.SafeString(kite.EmbedFs.ReadFile(path))
		_, err = fsutil.PutContents(dstFile, text, fsutil.FsCWTFlags)
		if err != nil {
			cliutil.Errorln("   ERROR")
			return err
		}
		cliutil.Successln("   OK")
	}

	return nil
}
