package main

import (
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/i18n"
	"github.com/gookit/kite/cmd"
	"github.com/gookit/slog"
)

var configFile string

func init() {
	config.AddDriver(yaml.Driver)
	err := config.LoadExists("kite.yaml")
	if err != nil {
		color.Error.Println("load config error:", err)
	}

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
	})

	langDir := "resource/language"
	slog.Println("load language files from:", langDir)
	i18n.Init(langDir, "zh-CN", map[string]string{
		"en":    "English",
		"zh-CN": "简体中文",
	})
}

func main() {
	app := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite Application"
		a.Description = "CLI tool application"
	})
	app.GOptsBinder = func(gf *gcli.Flags) {
		gf.StrOpt(
			&configFile,
			"config",
			"c",
			"kite.ini",
			"the YAML config file for kite",
		)
	}
	app.On(gcli.EvtAppPrepareAfter, func(_ ...interface{}) {
		if configFile == "" {
			return
		}

		slog.Printf("load custom config file %s", configFile)
		err := config.LoadFiles(configFile)
		if err != nil {
			color.Error.Println("load user config error:", err)
		}
	})

	cmd.AddCommands(app)

	app.Run()
}
