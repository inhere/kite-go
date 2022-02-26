package self

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/sysutil"
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/pkg/conf"
)

// KiteManage manage kite self
var KiteManage = &gcli.Command{
	Name: "self",
	Desc: "provide commands for manage kite self",
	Subs: []*gcli.Command{
		InitKite,
		KiteInfo,
		UpdateSelf,
	},
}

var KiteInfo = &gcli.Command{
	Name: "info",
	Desc: "show the kite tool information",
	Func: func(c *gcli.Command, args []string) error {
		show.AList("information", map[string]interface{}{
			"bin Dir":      c.BinDir(),
			"work Dir":     c.WorkDir(),
			"home Dir":     sysutil.HomeDir(),
			"loaded files": conf.C().LoadedFiles(),
			"language":     "TODO",
			"version":      kite.Version,
			"build date":   kite.BuildDate,
			"go version":   kite.GoVersion,
			// "i18n files": i18n.Default().LoadFile(),
		}, nil)

		return nil
	},
}

var UpdateSelf = &gcli.Command{
	Name:    "update",
	Aliases: []string{"updateself", "selfup", "upself", "up"},
	Desc:    "update {$binName} to latest from github repository",
}

var InitKite = &gcli.Command{
	Name: "init",
	Desc: "init kite env, will create config to use home dir",
}
