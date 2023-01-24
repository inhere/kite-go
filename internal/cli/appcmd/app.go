package appcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
)

// ManageCmd manage kite self
var ManageCmd = &gcli.Command{
	Name:    "app",
	Aliases: []string{"self"},
	Desc:    "provide commands for manage kite self",
	Subs: []*gcli.Command{
		AppCheckCmd,
		KiteInitCmd,
		KiteInfoCmd,
		UpdateSelfCmd,
		KiteConfCmd,
	},
}

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
			// "i18n files": i18n.Default().LoadFile(),
		}, nil)

		return nil
	},
}

// UpdateSelfCmd command
var UpdateSelfCmd = &gcli.Command{
	Name:    "update",
	Aliases: []string{"update-self", "up-self", "up"},
	Desc:    "update {$binName} to latest from github repository",
}
