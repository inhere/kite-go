package app

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/pkg/gitx"
	"github.com/inhere/kite-go/pkg/gitx/github"
	"github.com/inhere/kite-go/pkg/gitx/gitlab"
	"github.com/inhere/kite-go/pkg/httptpl"
	"github.com/inhere/kite-go/pkg/kiteext"
	"github.com/inhere/kite-go/pkg/kscript"
	"github.com/inhere/kite-go/pkg/lcproxy"
	"github.com/inhere/kite-go/pkg/quickjump"
)

const (
	ObjCli = "cli"
	ObjRux = "rux"

	ObjConf = "config"
	ObjLog = "logger" // console logger

	ObjPlugin = "plugin"
	ObjScript = "script"

	ObjGit  = "git"
	ObjGlab = "gitlab"
	ObjGhub = "github"
)

var (
	// L kite console logger
	L *slog.Logger
	// SL server logger
	SL *slog.Logger
)

var (
	Lcp *lcproxy.LocalProxy
	Exts *kiteext.ExtManager
	// AlsRun  *kiteext.KiteAliasRun

	Scripts *kscript.Runner
	Plugins *kiteext.PluginRunner

	QJump *quickjump.QuickJump

	// PathMap data
	PathMap *kiteext.PathMap
	OpenMap maputil.Aliases

	// Vars global vars.
	Vars *kiteext.VarMap
	HTpl *httptpl.Manager
)

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

// Ghub config get
func Ghub() *github.GitHub {
	return Get[*github.GitHub](ObjGhub)
}

// Gitx config get
func Gitx() *gitx.Config {
	return Get[*gitx.Config](ObjGit)
}
