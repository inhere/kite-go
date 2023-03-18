package pkgutil

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
)

// NewConfig box instance
func NewConfig() *config.Config {
	return config.
		NewWithOptions("kite", config.ParseEnv, config.ParseDefault, config.WithTagName("json")).
		WithDriver(yaml.Driver, JSON5Driver, ini.Driver)
}
