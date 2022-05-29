package bootstrap

import (
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/i18n"
	"github.com/gookit/slog"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/internal/command"
)

// Boot app
func Boot(app *app.KiteApp) error {
	LogBoot(app)

	slog.Info("bootstrap the kite application")

	if app.IsDebug() {
		dump.P(app.Cfg().Data())
	}

	// lang
	langDir := "resource/language"
	if fsutil.IsDir(langDir) {
		slog.Println("load language files from:", langDir)
		i18n.Init(langDir, "zh-CN", map[string]string{
			"en":    "English",
			"zh-CN": "简体中文",
		})
	}

	// load commands
	command.Boot(app)

	return nil
}
