package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/initlog"
)

// BootLogger handle
func BootLogger(ka *app.KiteApp) error {
	// output log to file
	logger := slog.New()
	logger.CallerFlag = slog.CallerFlagFull
	app.L = logger
	app.Add(app.ObjLog, logger)

	confMap := app.Cfg().SubDataMap(app.ObjLog)
	if len(confMap) == 0 {
		initlog.L.Info("skip init the kite logger, not found config")
		return nil
	}

	logCfg := handler.NewConfig(
		handler.WithRotateTime(rotatefile.EveryDay),
		handler.WithRotateMode(rotatefile.ModeCreate),
		handler.WithLogfile(ka.PathResolve(confMap.Str("logfile"))),
		handler.WithLevelName(confMap.Str("level")),
		handler.WithLevelModeString(confMap.Str("level_mode")),
		handler.WithBuffSize(confMap.Int("buffer_size")),
		handler.WithBackupNum(confMap.Uint("backup_num")),
	)
	// clean log files when close
	logCfg.CleanOnClose = true

	initlog.L.Infof("init the kite logger, rotate_mode:%s logfile: %s", logCfg.RotateMode, logCfg.Logfile)

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
