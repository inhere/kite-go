package appconf

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/structs"
	"github.com/inherelab/kite/internal/appconst"
	"github.com/inherelab/kite/internal/apputil"
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
	TmpDir    string   `json:"tmp_dir"`
	CacheDir  string   `json:"cache_dir"`
	HomeDir   string   `json:"home_dir"`
	WorkDir   string   `json:"work_dir"`
	PluginDir string   `json:"plugin_dir"`
	ConfFiles []string `json:"conf_files"`
}

// Init config
func Init() error {
	confFile := apputil.FindConfFile()
	if confFile != "" {
		err := confObj.LoadFiles(confFile)
		if err != nil {
			return err
		}

		// map config
		err = confObj.MapOnExists(appconst.ConfKeyApp, Conf)
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

var Aliases = &structs.Aliases{}
