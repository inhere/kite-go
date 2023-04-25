package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

// KiteInfoCmd instance
var KiteInfoCmd = &gcli.Command{
	Name: "info",
	Desc: "show the kite tool information",
	Func: func(c *gcli.Command, args []string) error {
		show.AList("information", map[string]interface{}{
			"user home dir": sysutil.UserHomeDir(),
			"app bin dir":   c.Ctx.BinDir(),
			"app bin file":  c.Ctx.BinFile(),
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
	keys   bool
	all    bool
}{}

// KiteConfCmd command
var KiteConfCmd = &gcli.Command{
	Name:    "config",
	Aliases: []string{"conf", "cfg"},
	Desc:    "display kite config information",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&confOpts.all, "all,a", "display all config data")
		c.BoolOpt2(&confOpts.keys, "keys", "display raw config data")
		c.StrOpt2(&confOpts.search, "search,s", "search top key by input keywords")

		c.AddArg("key", "show config for the input key")
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

		if confOpts.all {
			c.Infoln("All Config Data:")
			dump.Clear(app.Cfg().Data())
			return nil
		}

		key := c.Arg("key").String()
		if key == "" {
			return errorx.Raw("please input key for show configuration")
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

var kpOpts = struct {
	list bool
}{}

// KitePathCmd command
var KitePathCmd = &gcli.Command{
	Name: "path",
	// Aliases: []string{"update-self", "up-self", "up"},
	Desc: "show the path info on app by input name",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&kpOpts.list, "list, all, a, l", "display all paths for the kite")
		c.AddArg("name", "special path name on the kite, allow: base, config, tmp")
	},
	Func: func(c *gcli.Command, args []string) error {
		if kpOpts.list {
			show.AList("Kite paths", app.App().Config)
			return nil
		}

		name := c.Arg("name").String()
		if name == "" {
			return errorx.Raw("please input name for show path")
		}

		var path = app.App().PathByName(name)
		if path == "" {
			return errorx.Rawf("not found path for %q", name)
		}

		fmt.Println(path)
		return nil
	},
}

// kaCmdOpts struct
type kaCmdOpts struct {
	List bool `flag:"list all app command aliases;;;l"`
	Name string
}

var kaOpts = &kaCmdOpts{}

// KiteAliasCmd command
var KiteAliasCmd = &gcli.Command{
	Name:    "alias",
	Aliases: []string{"aliases", "cmd-alias"},
	Desc:    "show custom command aliases in app(config:aliases)",
	Config: func(c *gcli.Command) {
		goutil.MustOK(c.UseSimpleRule().FromStruct(kaOpts))
		c.AddArg("name", "get real-name of the input alias").WithAfterFn(func(a *gflag.CliArg) error {
			kaOpts.Name = a.String()
			return nil
		})
	},
	Func: func(c *gcli.Command, _ []string) error {
		if kaOpts.List {
			show.AList("Command aliases", cmdbiz.Kas)
			return nil
		}

		if kaOpts.Name != "" {
			fmt.Println(cmdbiz.Kas.ResolveAlias(kaOpts.Name))
			return nil
		}
		return errorx.New("please input alias for get command")
	},
}
