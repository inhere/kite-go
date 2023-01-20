package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/initlog"
)

// BootLogger handle
func BootLogger(ka *app.KiteApp) error {
	confMap := app.Cfg().SubDataMap("logger")
	logCfg := handler.NewConfig(
		handler.WithLogLevel(slog.LevelByName(confMap.Str("level"))),
		handler.WithLogfile(ka.PathResolve(confMap.Str("logfile"))),
		handler.WithLevelMode(uint8(confMap.Uint("level_mode"))),
		handler.WithBuffSize(confMap.Int("buffer_size")),
	)

	initlog.L.Infof("init the kite logger, logfile: %s", logCfg.Logfile)

	// output log to file
	logger := slog.NewWithConfig(func(l *slog.Logger) {

	})

	h1, err := logCfg.CreateHandler()
	if err != nil {
		return err
	}

	logger.AddHandlers(h1)

	// flush on end
	ka.OnShutdown(func() {
		logger.MustFlush()
	})

	app.Add(app.ObjLog, logger)
	return nil
}
