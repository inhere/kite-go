package app

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/i18n"
	"github.com/gookit/slog"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/app/appconf"
)

// LogConf struct
type LogConf struct {
	LogDir  string `json:"log_dir"`
	LogFile string `json:"log_file"`
}

// Boot app
func Boot(cli *gcli.App) {

	slog.Info("bootstrap kite runtime environment")

	err := appconf.Init()
	if err != nil {
		panic(err)
	}

	if IsDebug() {
		dump.P(appconf.C().Data())
	}

	Info.Version = kite.Version
	Info.GoVersion = kite.GoVersion

	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.WarnLevel

		f := logger.Formatter.(*slog.TextFormatter)
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}}\n")
	})
	// TODO output log to file

	// lang
	langDir := "resource/language"
	if fsutil.IsDir(langDir) {
		slog.Println("load language files from:", langDir)
		i18n.Init(langDir, "zh-CN", map[string]string{
			"en":    "English",
			"zh-CN": "简体中文",
		})
	}
}
