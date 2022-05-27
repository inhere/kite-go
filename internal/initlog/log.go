package initlog

import (
	"github.com/gookit/slog"
)

// L logger for init app
var L = slog.NewStdLogger().Config(func(sl *slog.SugaredLogger) {
	sl.Level = slog.DebugLevel

	f := sl.Formatter.(*slog.TextFormatter)
	if sl.Level >= slog.DebugLevel {
		f.SetTemplate("[{{level}}] {{caller}} {{message}} {{data}}\n")
	} else {
		f.SetTemplate("[{{level}}] {{message}} {{data}}\n")
	}
})
