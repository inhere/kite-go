package initlog

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/slog"
)

// DebugTpl for logger
const DebugTpl = "Kite [{{level}}] {{caller}} {{message}} {{data}}\n"

// L logger for init app
var L *slog.SugaredLogger

// Init logger
func Init(envLvName string) error {
	L = slog.NewStdLogger(func(sl *slog.SugaredLogger) {
		// sl.CallerFlag = slog.CallerFlagFull
		sl.CallerFlag = slog.CallerFlagFnlFcn
	})

	// init level
	lv := slog.LevelByName(envutil.Getenv(envLvName, "debug"))
	SetLevel(lv)

	L.Debug("the initlog create and init complete, level", L.Level.Name())
	return nil
}

// SetLevel for logger
func SetLevel(lv slog.Level) {
	L.Level = lv

	f := L.Formatter.(*slog.TextFormatter)
	if lv >= slog.DebugLevel {
		f.SetTemplate(DebugTpl)
	} else {
		f.SetTemplate("Kite [{{level}}] {{message}} {{data}}\n")
	}
}
