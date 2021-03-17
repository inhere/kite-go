package conf

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

var configObj = config.NewWith("kite", func(c *config.Config) {
	c.AddDriver(yamlv3.Driver)
	c.Options().TagName = "json"
})

// Conf for kite
var Conf = &Config{
	// LogFile: "./tmp/kite.log",
}

// Config struct
type Config struct {
	LogDir string `json:"log_dir"`
	LogFile string `json:"log_file"`
	TmpDir string `json:"tmp_dir"`
	CacheDir string `json:"cache_dir"`
	HomeDir string `json:"home_dir"`
	WorkDir string
}

// Obj get the config.Config
func Obj() *config.Config {
	return configObj
}
