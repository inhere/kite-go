package app

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/slog"
)

// Info for kite app
var Info = &struct {
	Branch    string
	Version   string
	PubDate   string
	Revision  string
	BuildDate string
	GoVersion string
}{
	Version: "1.0.0",
	PubDate: "2021-02-14 13:14",
}

var (
	KiteMode = envutil.Getenv("KITE_DEBUG", slog.DebugLevel.LowerName())
)

// IsDebug mode
func IsDebug() bool {
	return slog.LevelByName(KiteMode) >= slog.DebugLevel
}
