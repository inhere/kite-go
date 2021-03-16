package boot

import (
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/i18n"
	"github.com/gookit/slog"
	"github.com/inherelab/kite/pkg/conf"
)

func Boot() {
	// config
	conf.Obj().AddDriver(yamlv3.Driver)

	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
	})

	// lang
	langDir := "resource/language"
	slog.Println("load language files from:", langDir)
	i18n.Init(langDir, "zh-CN", map[string]string{
		"en":    "English",
		"zh-CN": "简体中文",
	})
}
