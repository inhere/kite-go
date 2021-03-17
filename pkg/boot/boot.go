package boot

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/i18n"
	"github.com/gookit/slog"
	"github.com/inherelab/kite/pkg/conf"
)

func Boot(app *gcli.App) {
	// config
	if conf.Obj().Exists("kite") {
		err := conf.Obj().MapStruct("kite", conf.Conf)
		if err != nil {
			color.Error.Println(err)
			return
		}
	}

	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
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
