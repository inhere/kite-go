package conf

import "github.com/gookit/config/v2"

var Conf = &Config{}

// Config struct
type Config struct {
	TmpDir string `json:"tmp_dir"`
	CacheDir string `json:"cache_dir"`
	HomeDir string `json:"home_dir"`
}

var configObj = config.NewWithOptions("kite", func(opt *config.Options) {
	opt.TagName = "json"
})

func Obj() *config.Config {
	return configObj
}
