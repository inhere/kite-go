package self

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/inherelab/kite/pkg/conf"
)

// KiteManage manage kite self
var KiteManage = &gcli.Command{
	Name: "self",
	Desc: "provide commands for manage kite self",
	Subs: []*gcli.Command{
		KiteInfo,
		UpdateSelf,
	},
}

var KiteInfo = &gcli.Command{
	Name: "info",
	Desc: "show the kite tool information",
	Func: func(c *gcli.Command, args []string) error {
		show.AList("information", map[string]interface{}{
			"binDir":  c.BinDir(),
			"workDir": c.WorkDir(),
			"loaded": conf.Obj().LoadedFiles(),
		}, nil)
		return nil
	},
}

var UpdateSelf = &gcli.Command{
	Name: "upself",
	Desc: "update {$binName} to latest from github repository",
}
