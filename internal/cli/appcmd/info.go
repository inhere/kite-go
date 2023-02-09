package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/internal/app"
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

// KiteObjectCmd command
var KiteObjectCmd = &gcli.Command{
	Name:    "object",
	Aliases: []string{"obj"},
	Desc:    "display service object config struct on kite",
	Config: func(c *gcli.Command) {
		c.AddArg("name", "show info for the object")
	},
	Func: func(c *gcli.Command, args []string) error {
		var data any
		key := c.Arg("name").String()
		switch key {
		case "app", "kite":
			data = app.App().Config
		case "git", "gitx":
			data = app.Gitx()
		case "glab", "gitlab":
			data = app.Glab()
		default:
			if app.Has(key) {
				data = app.GetAny(key)
			}
		}

		if data == nil {
			return c.NewErrf("not found object for %q", key)
		}

		c.Warnln("Object info for", key)
		dump.Clear(data)
		return nil
	},
}

var confOpts = struct {
	search string
	raw    bool
	keys   bool
}{}

// KiteConfCmd command
var KiteConfCmd = &gcli.Command{
	Name:    "config",
	Aliases: []string{"conf", "cfg"},
	Desc:    "display kite config information",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&confOpts.raw, "raw", "r", false, "display raw config data")
		c.BoolOpt2(&confOpts.keys, "keys", "display raw config data")
		c.StrOpt2(&confOpts.search, "search,s", "search top key by input keywords")

		c.AddArg("key", "show config for the key")
	},
	Func: func(c *gcli.Command, args []string) error {
		if confOpts.keys {
			names := make([]string, 16)
			for name := range app.Cfg().Data() {
				names = append(names, name)
			}

			c.Infoln("All Config Keys:")
			dump.Clear(names)
			return nil
		}

		key := c.Arg("key").String()
		if key == "" {
			c.Infoln("All Config Data:")
			dump.Clear(app.Cfg().Data())
			return nil
		}

		switch key {
		case "git":
			key = app.ObjGit
		case "glab":
			key = app.ObjGlab
		case "hub", "ghub":
			key = app.ObjGhub
		}

		data, ok := app.Cfg().GetValue(key)
		if !ok {
			return c.NewErrf("not found config for key: %s", key)
		}

		c.Infoln("Config for key:", key)
		dump.Clear(data)
		return nil
	},
}

// PathAliasCmd command
var PathAliasCmd = &gcli.Command{
	Name:    "pathmap",
	Aliases: []string{"path-alias"},
	Desc:    "custom path aliases mapping in kite",
	Func: func(c *gcli.Command, args []string) error {
		return errorx.New("todo")
	},
}

// KiteAliasCmd command
var KiteAliasCmd = &gcli.Command{
	Name:    "alias",
	Aliases: []string{"aliases"},
	Desc:    "display custom command aliases for kite",
	Func: func(c *gcli.Command, args []string) error {
		return errorx.New("todo")
	},
}

// CommandMapCmd command
var CommandMapCmd = &gcli.Command{
	Name:    "cmd-map",
	Aliases: []string{"cmdmap"},
	Desc:    "display all console commands info for kite",
	Func: func(c *gcli.Command, args []string) error {
		return errorx.New("todo")
	},
}
