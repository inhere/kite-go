package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/github"
	"github.com/inhere/kite/pkg/gitx/gitlab"
)

const (
	ObjCli = "cli"
	ObjRux = "rux"

	ObjConf = "config"
	ObjLog  = "logger"

	ObjPlugin = "plugin"
	ObjScript = "script"

	ObjGit  = "git"
	ObjGlab = "gitlab"
	ObjGhub = "github"
)

// L kite logger
var L *slog.Logger

// SL server logger
var SL *slog.Logger

// Cfg get the config object
func Cfg() *config.Config {
	return Get[*config.Config](ObjConf)
}

// Rux get the web app
func Rux() *rux.Router {
	return Get[*rux.Router](ObjRux)
}

// Cli get the cli app
func Cli() *gcli.App {
	return Get[*gcli.App](ObjCli)
}

// Log get
func Log() *slog.Logger {
	return Get[*slog.Logger](ObjLog)
}

// Glab get
func Glab() *gitlab.GitLab {
	return Get[*gitlab.GitLab](ObjGlab)
}

// Ghub get
func Ghub() *github.GitHub {
	return Get[*github.GitHub](ObjGhub)
}

// Gitx get
func Gitx() *gitx.Config {
	return Get[*gitx.Config](ObjGit)
}
