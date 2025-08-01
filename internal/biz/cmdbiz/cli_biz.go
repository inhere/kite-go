package cmdbiz

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/lcproxy"
)

// CommonOpts some common vars struct
type CommonOpts struct {
	Proxy   bool
	DryRun  bool `flag:"name=dry-run;desc=run workflow, but dont real execute command;shorts=dry"`
	Confirm bool
	Workdir string
}

// BindCommonFlags for some git commands
func (co *CommonOpts) BindCommonFlags(c *gcli.Command) {
	co.BindWorkdirDryRun(c)
	co.BindProxyConfirm(c)
}

// BindCommonFlags2 for some git commands -w: 不会自动改变工作目录
func (co *CommonOpts) BindCommonFlags2(c *gcli.Command) {
	co.BindWorkdirDryRun2(c)
	co.BindProxyConfirm(c)
}

// BindProxyConfirm flags for some cli commands
func (co *CommonOpts) BindProxyConfirm(c *gcli.Command) {
	c.BoolOpt2(&co.Proxy, "proxy,P", "manual enable set proxy ENV config(config:local_proxy)", gflag.WithValidator(func(val string) error {
		if strutil.SafeBool(val) {
			app.Lcp.Apply(func(lp *lcproxy.LocalProxy) {
				c.Infoln("TIP: enabled to set proxy ENV vars, will set", lp.EnvKeys())
				dump.NoLoc(lp)
			})
		}
		return nil
	}))
	c.BoolOpt2(&co.Confirm, "confirm,C", "confirm ask before executing command")
}

// BindWorkdirDryRun flags for some cli commands
func (co *CommonOpts) BindWorkdirDryRun(c *gcli.Command) {
	c.BoolOpt(&co.DryRun, "dry-run", "dry", false, "run workflow, but dont real execute command")
	c.StrOpt2(&co.Workdir, "workdir, w", "the command workdir path, default is current dir", gflag.WithHandler(func(val string) error {
		realPath := fsutil.Realpath(val)
		slog.Debugf("will change current workdir to: %s", realPath)
		return c.ChWorkDir(realPath)
	}))
}

// BindWorkdirDryRun2 flags for some cli commands 不会自动改变工作目录
func (co *CommonOpts) BindWorkdirDryRun2(c *gcli.Command) {
	c.BoolOpt(&co.DryRun, "dry-run", "dry", false, "run workflow, but dont real execute command")
	c.StrOpt2(&co.Workdir, "workdir, w", "the command workdir path, default is current dir")
}
