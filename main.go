package main

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/ini/v2"
	"github.com/gookit/kite/cmd/swagger"
	"github.com/gookit/slog"
)

var configFile string

func init() {
	err := ini.LoadExists("kite.ini")
	if err != nil {
		color.Error.Println("load config error:", err)
	}

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
	})
}

func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "kite Application"
		a.Description = "CLI tool application for rux"
	})
	app.GOptsBinder = func(gf *gcli.Flags) {
		gf.StrOpt(&configFile, "config", "c", "kite.ini", "the INI config file for kite")
	}

	loadCommands(app)

	app.Run()
}

func loadCommands(app *gcli.App) {
	app.AddCommand(swagger.GenCode)
	app.AddCommand(swagger.DocBrowse)
	app.AddCommand(swagger.DocGen)
}
