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
//
// 注意: 这些字段统一通过下面的 Bind* 方法绑定为命令选项, 不要加 flag tag。
// 否则当命令匿名内嵌 CommonOpts 并用 FromStruct 时, 会与 Bind 方法重复注册同名
// 选项而 panic(gcli v3.5+ 的 FromStruct 会自动展开匿名内嵌结构体的 tag)。
type CommonOpts struct {
	Proxy   bool
	DryRun  bool
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
