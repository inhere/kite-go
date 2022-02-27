package conf

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/sysutil"
)

var confObj = config.NewWith("kite", func(c *config.Config) {
	c.AddDriver(yamlv3.Driver)
	c.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "json"
	})
})

// Conf for kite
var Conf = &Config{
	// LogFile: "./tmp/kite.log",
}

// Config struct
type Config struct {
	TmpDir    string `json:"tmp_dir"`
	CacheDir  string `json:"cache_dir"`
	HomeDir   string `json:"home_dir"`
	WorkDir   string
	PluginDir string   `json:"plugin_dir"`
	ConfFiles []string `json:"conf_files"`
}

// Init config
func Init() error {
	confFile := findConfFile()
	if confFile != "" {
		err := confObj.LoadFiles(confFile)
		if err != nil {
			return err
		}

		// map config
		err = MapOnExists("kite", Conf)
		if err != nil {
			return err
		}
	}

	return nil
}

// C get the config.Config
func C() *config.Config {
	return confObj
}

// MapOnExists by key name.
func MapOnExists(key string, ptr interface{}) error {
	if confObj.Exists(key) {
		err := confObj.MapStruct(key, Conf)
		if err != nil {
			return err
		}
	}

	return nil
}

var Aliases = &structs.Aliases{}

func findConfFile() string {
	confFile := envutil.Getenv("KITE_CONF_FILE", sysutil.UserDir(".kite/kite.yml"))
	if fsutil.IsFile(confFile) {
		return confFile
	}

	confFile = cliutil.Workdir() + "/kite.yml"
	if fsutil.IsFile(confFile) {
		return confFile
	}

	confFile = cliutil.BinDir() + "/kite.yml"
	if fsutil.IsFile(confFile) {
		return confFile
	}

	return ""
}
