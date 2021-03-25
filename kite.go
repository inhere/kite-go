package kite

// Info for kite app
var Info = &struct {
	Version string
	PubDate string
}{
	Version: "1.0.0",
	PubDate: "2021-02-14 13:14",
}

var (
	Version = "0.0.0"
	PubDate = "2021-02-14 13:14"
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
