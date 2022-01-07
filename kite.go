package kite

var (
	Version   = "0.0.0"
	PubDate   = "2021-02-14 13:14"
	Branch    string
	Revision  string
	BuildDate string
	GoVersion string
)

// Env names
const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvPre   = "pre"
	EnvProd  = "prod"
)

var (
	Timezone   = "PRC"
	DateFormat = "2006-01-02 15:04:05"
)
