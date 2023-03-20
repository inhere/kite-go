package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/initlog"
)

// BootLogger handle
func BootLogger(ka *app.KiteApp) error {
	// output log to file
	logger := slog.New()
	logger.CallerFlag = slog.CallerFlagFull
	app.Add(app.ObjLog, logger)

	confMap := app.Cfg().SubDataMap(app.ObjLog)
	if len(confMap) == 0 {
		initlog.L.Info("skip init the kite logger, not found config")
		return nil
	}

	logCfg := handler.NewConfig(
		handler.WithLogLevel(slog.LevelByName(confMap.Str("level"))),
		handler.WithLogfile(ka.PathResolve(confMap.Str("logfile"))),
		handler.WithLevelMode(uint8(confMap.Uint("level_mode"))),
		handler.WithBuffSize(confMap.Int("buffer_size")),
	)

	initlog.L.Infof("init the kite logger, logfile: %s", logCfg.Logfile)

	h1, err := logCfg.CreateHandler()
	if err != nil {
		return err
	}

	logger.AddHandlers(h1)

	// flush on end
	ka.OnShutdown(func() {
		logger.MustFlush()
	})

	return nil
}

func BootSrvLogger(ka *app.KiteApp) error {
	return nil
}
