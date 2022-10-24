package initlog

import (
	"github.com/gookit/slog"
)

// L logger for init app
var L = slog.NewStdLogger().Config(func(sl *slog.SugaredLogger) {
	sl.Level = slog.DebugLevel
	sl.CallerFlag = slog.CallerFlagFull

	f := sl.Formatter.(*slog.TextFormatter)
	if sl.Level >= slog.DebugLevel {
		f.SetTemplate("Kite [{{level}}] {{caller}} {{message}} {{data}}\n")
	} else {
		f.SetTemplate("Kite [{{level}}] {{message}} {{data}}\n")
	}
})
