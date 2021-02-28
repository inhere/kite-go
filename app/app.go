package app

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

const (
	GhContentHost = "https://raw.githubusercontent.com"
	// eg: https://raw.githubusercontent.com/gookit/slog/master/README.md
	GhContentURL = "https://raw.githubusercontent.com/%s/%s/%s"
)

type Option struct {
	CacheDir string
}
