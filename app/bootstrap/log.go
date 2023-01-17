package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/inhere/kite/app"
)

func BootLogger(ka *app.KiteApp) error {

	// slog
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.WarnLevel

		f := logger.Formatter.(*slog.TextFormatter)
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}}\n")
	})

	// TODO output log to file

	return nil
}
