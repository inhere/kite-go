package main

import (
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/slog"
	"github.com/inherelab/kite/cmd"
	"github.com/inherelab/kite/pkg/conf"
)

var confFile string

// dev run:
//	go run ./cmd/kit
func init() {


}

func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"
	})
	app.GOptsBinder = func(gf *gcli.Flags) {
		gf.StrOpt(&confFile,
			"config",
			"c",
			"kite.yaml",
			"the YAML config file for kite",
		)
	}
	app.On(gcli.EvtGOptionsParsed, func(_ ...interface{}) {
		if confFile == "" {
			return
		}

		slog.Printf("load custom config file %s", confFile)
		err := config.LoadFiles(confFile)
		if err != nil {
			color.Error.Println("load user config error:", err)
		}
	})

	boot()

	cmd.Register(app)
	app.Run(nil)
}

func boot() {
	err := conf.Obj().MapStruct("kite", conf.Conf)
	if err != nil {
		color.Error.Println(err)
		return
	}

}
