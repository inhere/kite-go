package initlog

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/slog"
)

// L logger for init app
var L *slog.SugaredLogger

// Init logger
func Init(envLvName string) error {
	L = slog.NewStdLogger(func(sl *slog.SugaredLogger) {
		sl.CallerFlag = slog.CallerFlagFnlFcn
		// sl.CallerFlag = slog.CallerFlagFull
		sl.Level = slog.LevelByName(envutil.Getenv(envLvName, "debug"))

		f := sl.Formatter.(*slog.TextFormatter)
		if sl.Level >= slog.DebugLevel {
			f.SetTemplate("Kite [{{level}}] {{caller}} {{message}} {{data}}\n")
		} else {
			f.SetTemplate("Kite [{{level}}] {{message}} {{data}}\n")
		}
	})

	L.Debug("the initlog create and init complete, level", L.Level.Name())
	return nil
}
