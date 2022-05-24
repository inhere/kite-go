package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/inherelab/kite/app"
)

func LogBoot(app *app.KiteApp) {
	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.WarnLevel

		f := logger.Formatter.(*slog.TextFormatter)
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}}\n")
	})

	// TODO output log to file

}
