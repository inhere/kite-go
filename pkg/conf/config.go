package conf

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/structs"
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
	LogDir   string `json:"log_dir"`
	LogFile  string `json:"log_file"`
	TmpDir   string `json:"tmp_dir"`
	CacheDir string `json:"cache_dir"`
	HomeDir  string `json:"home_dir"`
	WorkDir  string
}

// C get the config.Config
func C() *config.Config {
	return confObj
}

var Aliases = &structs.Aliases{}
