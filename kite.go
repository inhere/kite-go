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

// AppConfig struct
//
// Gen by:
//   kite go gen st -s @c -t json --name AppConfig
type AppConfig struct {
	// BaseDir base dir
	BaseDir string `json:"base_dir"`
	// TmpDir tmp dir
	TmpDir string `json:"tmp_dir"`
	// CacheDir cache dir
	CacheDir string `json:"cache_dir"`
	// ConfigDir config dir
	ConfigDir string `json:"config_dir"`
	// ResourceDir resource dir
	ResourceDir string `json:"resource_dir"`
	// IncludeConfig include config
	IncludeConfig string `json:"include_config"`
}

// App struct
type App struct {
	AppConfig
}
