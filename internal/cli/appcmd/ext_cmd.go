package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/gookit/slog"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kiteext"
)

var extRunOpts = struct {
	EnvMap gflag.KVString `flag:"name=env;desc=set env vars for run, allow input multi;short=e"`
	// DryRun setting
	DryRun  bool   `flag:"name=dry-run;desc=Set whether it is actually execute ext;shorts=dry"`
	Verbose bool   `flag:"name=verbose;desc=Set whether to print verbose information;shorts=v"`
	Workdir string `flag:"desc=set workdir for run ext;shorts=w"`
}{}

// NewAppExtCmd 创建应用扩展管理命令
func NewAppExtCmd() *gcli.Command {
	return &gcli.Command{
		Name: "ext",
		Desc: "Kite extensions manage command",
		Help: `Quick Run:
  {$binName} <ext> [args for ext...]
  {$binWithCmd} <ext> [args for ext...]
`,
		Subs: []*gcli.Command{
			AppExtListCmd,
			AppExtAddCmd,
			AppExtUpdateCmd,
			AppExtRunCmd,
		},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&extRunOpts)
			c.On(events.OnCmdSubNotFound, SubCmdNotFound)
		},
	}
}

// AppExtRunCmd run ext commands
var AppExtRunCmd = &gcli.Command{
	Name: "run",
	Desc: "run extension commands with args",
	Help: `Example: kite app ext run [options] <ext> [args...]`,
	Config: func(c *gcli.Command) {
		c.MustFromStruct(&extRunOpts)
		c.AddArg("name", "extension name for run", true)
	},
	Func: func(c *gcli.Command, args []string) error {
		name := c.Arg("name").String()
		return handleExtRun(name, &kiteext.RunCtx{
			Dir:  extRunOpts.Workdir,
			Dry:  extRunOpts.DryRun,
			Env:  extRunOpts.EnvMap.Data(),
			Args: args,
		})
	},
}

func handleExtRun(name string, rc *kiteext.RunCtx) error {
	return app.Exts.Run(name, rc)
}

func SubCmdNotFound(ctx *gcli.HookCtx) (stop bool) {
	name := ctx.Str("name")
	args := ctx.Strings("args")
	slog.Debugf("ext.SubCmdNotFound: %s, args: %v, try run as ext name", name, args)

	err := handleExtRun(name, &kiteext.RunCtx{
		Dir:  extRunOpts.Workdir,
		Dry:  extRunOpts.DryRun,
		Env:  extRunOpts.EnvMap.Data(),
		Args: args,
	})
	if err == nil {
		return true
	}

	// ctx.WithErr(err)
	ccolor.Errorln("run ext error", err)
	return
}

var (
	AppExtListCmd = &gcli.Command{
		Name:    "list",
		Desc:    "list all installed extensions",
		Aliases: []string{"ls"},
		Config: func(c *gcli.Command) {
			c.AddArg("keywords", "keywords for search extensions", false, true)
		},
		Func: func(c *gcli.Command, args []string) error {
			ccolor.Cyanln("Extensions metafile: ", app.Exts.Metafile)
			ccolor.Cyanln("Installed extensions:")
			keywords := c.Arg("keywords").Array()

			for _, ext := range app.Exts.Exts() {
				if len(keywords) > 0 {
					if !strutil.ContainsOne(ext.Name, keywords) && !strutil.ContainsOne(ext.Desc, keywords) {
						continue
					}
				}

				version := strutil.OrCond(ext.Version != "", " (v"+ext.Version+")", "")
				ccolor.Printf("  <green>%s</>  %s%s\n", ext.Name, ext.Desc, version)
				fmt.Printf("    - path: %s\n", ext.OsPath())
			}

			return nil
		},
	}
)

// add new extension
var (
	extSetOpts = struct {
		Path     string `flag:"desc=set extension file path;shorts=p"`
		Desc     string `flag:"desc=set extension description;shorts=d"`
		Aliases  string `flag:"desc=set extension aliases, multi by comma;shorts=alias"`
		Author   string `flag:"desc=set extension author;shorts=user"`
		Version  string `flag:"desc=set extension version;shorts=v,ver"`
		BinName  string `flag:"desc=set extension binary name for search;shorts=bin"`
		Homepage string `flag:"desc=set extension homepage URL;shorts=home"`
		// 交互模式设置信息
		Interactive bool `flag:"desc=interactive mode;shorts=i"`
	}{}

	// AppExtAddCmd 添加应用扩展添加命令
	AppExtAddCmd = &gcli.Command{
		Name: "add",
		Desc: "Add a new kite extension",
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&extSetOpts)
			c.AddArg("name", "extension name", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()
			ext := kiteext.NewExt(name, extSetOpts.Desc, extSetOpts.Path)
			ext.SetBinName(extSetOpts.BinName)
			ext.BeforeSave = func(ext *kiteext.KiteExt) bool {
				show.AList("Extension Info", ext)
				return interact.Confirm("Do you want to save this extension?")
			}

			// TODO set more info

			return app.Exts.Add(ext)
		},
	}
)

// update kite extension
var (
	// AppExtUpdateCmd update kite extension
	AppExtUpdateCmd = &gcli.Command{
		Name: "update",
		Desc: "Update kite extension",
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&extSetOpts)
			c.AddArg("name", "extension name", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			name := c.Arg("name").String()
			ext, ok := app.Exts.Ext(name)
			if !ok {
				return fmt.Errorf("extension %s not found", name)
			}

			// set info
			ext.SetPath(extSetOpts.Path)
			ext.SetBinName(extSetOpts.BinName)
			ext.SetAliases(extSetOpts.Aliases)
			if extSetOpts.Desc != "" {
				ext.Desc = extSetOpts.Desc
			}
			if extSetOpts.Author != "" {
				ext.Author = extSetOpts.Author
			}
			if extSetOpts.Homepage != "" {
				ext.Homepage = extSetOpts.Homepage
			}
			if extSetOpts.Version != "" {
				ext.Version = extSetOpts.Version
			}

			ext.BeforeSave = func(ext *kiteext.KiteExt) bool {
				show.AList("Extension Info", ext)
				return interact.Confirm("Do you want to save this extension?")
			}
			return app.Exts.Update(ext)
		},
	}
)
