package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/inhere/kite/app"
)

var confOpts = struct {
	key string
	raw bool
}{}

// KiteConfCmd command
var KiteConfCmd = &gcli.Command{
	Name:    "config",
	Aliases: []string{"conf", "cfg"},
	Desc:    "display kite config information",
	Config: func(c *gcli.Command) {
		c.BoolOpt(&confOpts.raw, "raw", "r", false, "display raw config data")
		c.StrOpt(&confOpts.key, "key", "k", "show config for the key")
	},
	Func: func(c *gcli.Command, args []string) error {
		key := confOpts.key
		if key == "" {
			c.Infoln("Config for kite:")
			dump.Clear(app.Cfg().Data())
			return nil
		}

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
	},
}
