package main

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/cmd"
	"github.com/inherelab/kite/pkg/conf"
)

var confFile string

// dev run:
//	go run ./bin/kit
//	go run ./bin/kit
func main() {
	cli := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"

		a.Version = kite.Version
	})
	cli.GOptsBinder = func(gfs *gcli.Flags) {
		gfs.StrOpt(&confFile,
			"config",
			"c",
			"kite.yaml",
			"the YAML config file for kite",
		)
	}
	cli.On(gcli.EvtGOptionsParsed, func(_ ...interface{}) bool {
		if confFile != "" {
			color.Infoln("load custom config file:", confFile)
			err := conf.Obj().LoadExists(confFile)
			if err != nil {
				color.Error.Println("load user config error:", err)
				return false
			}
		}

		// boot kite
		color.Infoln("bootstrap kite runtime environment")
		app.Boot(cli)
		return false
	})

	// load commands
	cmd.Boot(cli)

	// do run
	cli.Run(nil)
}
