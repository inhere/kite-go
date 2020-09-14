package main

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/ini/v2"
	"github.com/gookit/ruxc/cmd"
)

var configFile string

func init() {
	err := ini.LoadExists("ruxc.ini")
	if err != nil {
		color.Error.Println("load config error:", err)
	}
}

func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Ruxc Application"
		a.Description = "CLI tool application for rux"

	})
	app.GOptsBinder = func(gf *gcli.Flags) {
		gf.StrOpt(&configFile, "config", "c", "ruxc.ini", "the INI config file for ruxc")
	}

	loadCommands(app)

	app.Run()
}

func loadCommands(app *gcli.App)  {
	app.AddCommand(cmd.Swag2code)
}
