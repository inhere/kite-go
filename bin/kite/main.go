package main

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/inherelab/kite/cmd"
	"github.com/inherelab/kite/pkg/boot"
	"github.com/inherelab/kite/pkg/conf"
)

var confFile string

// dev run:
//	go run ./bin/kit
//	go run ./bin/kit
func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"
	})
	app.GOptsBinder = func(gfs *gcli.Flags) {
		gfs.StrOpt(&confFile,
			"config",
			"c",
			"kite.yaml",
			"the YAML config file for kite",
		)
	}
	app.On(gcli.EvtGOptionsParsed, func(_ ...interface{}) {
		if confFile != "" {
			color.Infoln("load custom config file:", confFile)
			err := conf.Obj().LoadExists(confFile)
			if err != nil {
				color.Error.Println("load user config error:", err)
				return
			}
		}

		// boot kite
		color.Infoln("bootstrap kite runtime environment")
		boot.Boot(app)
	})

	// load commands
	cmd.Register(app)

	// do run
	app.Run(nil)
}
