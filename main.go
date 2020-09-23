package main

import (
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/i18n"
	"github.com/gookit/kite/cmd"
	"github.com/gookit/slog"
)

var configFile string

func init() {
	err := config.LoadExists("kite.yaml")
	if err != nil {
		color.Error.Println("load config error:", err)
	}

	i18n.Init("resource/language", "zh-CN", map[string]string{
		"en":"English",
		"zh-CN":"简体中文",
	})

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
	})
}

func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite Application"
		a.Description = "CLI tool application"
	})
	app.GOptsBinder = func(gf *gcli.Flags) {
		gf.StrOpt(&configFile,
			"config",
			"c",
			"kite.ini",
			"the YAML config file for kite",
			)
	}

	cmd.AddCommands(app)

	app.Run()
}
