package initlog

import (
	"github.com/gookit/slog"
)

var L = slog.NewStdLogger().Config(func(sl *slog.SugaredLogger) {
	sl.Level = slog.DebugLevel

	f := sl.Formatter.(*slog.TextFormatter)
	f.SetTemplate("[{{level}}] {{message}} {{data}}\n")
})
