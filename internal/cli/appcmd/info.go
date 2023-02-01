package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
)

// KiteInfoCmd instance
var KiteInfoCmd = &gcli.Command{
	Name: "info",
	Desc: "show the kite tool information",
	Func: func(c *gcli.Command, args []string) error {
		show.AList("information", map[string]interface{}{
			"user home dir": sysutil.UserHomeDir(),
			"kite bin dir":  c.Ctx.BinDir(),
			"user data dir": app.App().BaseDir,
			"work dir":      c.Ctx.WorkDir(),
			"dotenv file":   app.App().DotenvFile(),
			"config files":  app.Cfg().LoadedFiles(),
			"language":      "TODO",
			"version":       kite.Version,
			"build date":    kite.BuildDate,
			"go version":    kite.GoVersion,
			"github repo":   kite.GithubRepo,
			// "i18n files": i18n.Default().LoadFile(),
		}, nil)

		return nil
	},
}

var kpOpts = struct {
	all bool
}{}

// KitePathCmd command
var KitePathCmd = &gcli.Command{
	Name: "path",
	// Aliases: []string{"update-self", "up-self", "up"},
	Desc: "show the path info on kite by input name",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&kpOpts.all, "all, a", "display all paths for the kite")
		c.AddArg("name", "special path name on the kite, allow: base, config, tmp")
	},
	Func: func(c *gcli.Command, args []string) error {
		if kpOpts.all {
			dump.Clear(app.App().Config)
			return nil
		}

		name := c.Arg("name").String()
		if name == "" {
			return errorx.Raw("Please input name for show path")
		}

		var path string
		switch name {
		case "base", "root":
			path = app.App().BaseDir
		case "cache", "caches":
			path = app.App().CacheDir
		case "conf", "config":
			path = app.App().ConfigDir
		case "data":
			path = app.App().DataDir
		case "tmp", "temp":
			path = app.App().TmpDir
		case "res", "resource":
			path = app.App().ResourceDir
		}

		if path == "" {
			return errorx.Rawf("Not found path for %q", name)
		}

		fmt.Println(path)
		return nil
	},
}
