package initlog

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/slog"
)

// EnvInitLogLevel key
const EnvInitLogLevel = "KITE_INIT_LOG"

// L logger for init app
var L = slog.NewStdLogger().Config(func(sl *slog.SugaredLogger) {
	sl.CallerFlag = slog.CallerFlagFnlFcn
	// sl.CallerFlag = slog.CallerFlagFull
	// sl.CallerSkip += 1

	// sl.Level = slog.DebugLevel
	sl.Level = slog.LevelByName(envutil.Getenv(EnvInitLogLevel, "debug"))

	f := sl.Formatter.(*slog.TextFormatter)
	if sl.Level >= slog.DebugLevel {
		f.SetTemplate("Kite [{{level}}] {{caller}} {{message}} {{data}}\n")
	} else {
		f.SetTemplate("Kite [{{level}}] {{message}} {{data}}\n")
	}
})
