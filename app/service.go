package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/rux"
)

// Cfg get the config object
func Cfg() *config.Config {
	return Get[*config.Config]("config")
}

// Rux get the web app
func Rux() *rux.Router {
	return Get[*rux.Router]("config")
}

// Cli get the cli app
func Cli() *gcli.App {
	return Get[*gcli.App]("cli")
}