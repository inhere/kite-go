package self

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
)

// KiteManage manage kite self
var KiteManage = &gcli.Command{
	Name:    "app",
	Aliases: []string{"self"},
	Desc:    "provide commands for manage kite self",
	Subs: []*gcli.Command{
		InitKite,
		KiteInfo,
		UpdateSelf,
		KiteConf,
	},
}

var KiteInfo = &gcli.Command{
	Name: "info",
	Desc: "show the kite tool information",
	Func: func(c *gcli.Command, args []string) error {
		show.AList("information", map[string]interface{}{
			"bin dir":      c.Ctx.BinDir(),
			"home dir":     sysutil.HomeDir(),
			"work dir":     c.Ctx.WorkDir(),
			"config files": app.Cfg().LoadedFiles(),
			"language":     "TODO",
			"version":      kite.Version,
			"build date":   kite.BuildDate,
			"go version":   kite.GoVersion,
			// "i18n files": i18n.Default().LoadFile(),
		}, nil)

		return nil
	},
}

// UpdateSelf command
var UpdateSelf = &gcli.Command{
	Name:    "update",
	Aliases: []string{"updateself", "selfup", "upself", "up"},
	Desc:    "update {$binName} to latest from github repository",
}

// InitKite command
var InitKite = &gcli.Command{
	Name: "init",
	Desc: "init kite env, will create config to use home dir",
}

var (
	confOpts = struct {
		key string
		raw bool
	}{}

	// KiteConf command
	KiteConf = &gcli.Command{
		Name:    "config",
		Aliases: []string{"conf", "cfg"},
		Desc:    "show application config information",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&confOpts.raw, "raw", "r", false, "display raw config data")
			c.StrOpt(&confOpts.key, "key", "k", "show config for the key")
		},
		Func: func(c *gcli.Command, args []string) error {
			if key := confOpts.key; key != "" {
				var data any
				if !confOpts.raw {
					switch key {
					case "app":
						data = app.App().Config
					}
				}

				if data == nil {
					var ok bool
					data, ok = app.Cfg().GetValue(key)
					if !ok {
						return c.NewErrf("not found config for key: %s", key)
					}
				}

				c.Infoln("Config for key:", key)
				dump.Clear(data)
				return nil
			}

			c.Infoln("Config for app:")
			dump.Clear(app.Cfg().Data())
			return nil
		},
	}
)
