package appconst

const (
	// EnvKiteVerbose level.
	EnvKiteVerbose = "KITE_VERBOSE"
	// EnvKiteConfig main config file env name
	EnvKiteConfig = "KITE_CONFIG"
	// KiteConfigName default main config filename
	KiteConfigName = "kite.yml"
	// KiteDefaultDataDir path
	KiteDefaultDataDir = "~/.kite"
	// KiteDefaultConfigFile path
	KiteDefaultConfigFile = "~/.kite/" + KiteConfigName

	// AppName for the application
	AppName = "kite"

	// ConfKeyApp name on config
	ConfKeyApp = "app"
)

var (
	Timezone   = "PRC"
	DateFormat = "2006-01-02 15:04:05"
)
