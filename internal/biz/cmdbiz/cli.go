package cmdbiz

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/lcproxy"
)

// CommonOpts some common vars struct
type CommonOpts struct {
	Proxy   bool
	DryRun  bool
	Confirm bool
	Workdir string
	GitHost string
}

// BindCommonFlags for some git commands
func (co *CommonOpts) BindCommonFlags(c *gcli.Command) {
	co.BindCommonFlags1(c)

	c.BoolOpt2(&co.Proxy, "proxy,P", "manual enable set proxy ENV config", gflag.WithValidator(func(val string) error {
		if strutil.QuietBool(val) {
			app.App().Lcp.Apply(func(lp *lcproxy.LocalProxy) {
				c.Infoln("TIP: enabled to set proxy ENV vars, will set", lcproxy.HttpKey, lcproxy.HttpsKey)
				dump.NoLoc(lp)
			})
		}
		return nil
	}))
	c.BoolOpt2(&co.Confirm, "confirm", "confirm ask before executing command")
}

// BindCommonFlags1 for some git commands
func (co *CommonOpts) BindCommonFlags1(c *gcli.Command) {
	c.BoolOpt(&co.DryRun, "dry-run", "dry", false, "run workflow, but dont real execute command")
	c.StrOpt(&co.Workdir, "workdir", "w", "", "the command workdir path, default is current dir")
}
