package appconst

const (
	// EnvInitLogLevel key
	EnvInitLogLevel = "KITE_INIT_LOG"
	// EnvKiteVerbose level.
	EnvKiteVerbose = "KITE_VERBOSE"
	// EnvKiteBaseDir for override the KiteDefaultBaseDir
	EnvKiteBaseDir = "KITE_BASE_DIR"
	// EnvKiteConfig main config file env name
	EnvKiteConfig = "KITE_CONFIG_FILE"
	EnvKiteDotEnv = "KITE_DOTENV_FILE"

	DotEnvFileName = ".env"
	// KiteConfigName default main config filename
	KiteConfigName = "kite.yml"
	// KiteDefaultBaseDir path for: config, tmp and more
	KiteDefaultBaseDir = "~/.kite-go"

	// AppName for the application
	AppName = "kite"

	// ConfKeyApp name on config
	ConfKeyApp = "app"
)

var (
	Timezone   = "PRC"
	DateFormat = "2006/01/02 15:04:05.000"
)

const (
	VarFormat = "{,}"
)
