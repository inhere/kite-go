package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/ccolor"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/pkg/kiteext"
)

// NewAppExtCmd 创建应用扩展管理命令
func NewAppExtCmd() *gcli.Command {
	return &gcli.Command{
		Name: "ext",
		Desc: "Kite extensions manage command",
		Subs: []*gcli.Command{
			AppExtListCmd,
			AppExtAddCmd,
		},
	}
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

var (
	extAddOpts = struct {
		Path string `flag:"desc=set extension description;shorts=bin"`
		Desc string `flag:"desc=set extension description;shorts=d"`
		// 交互模式设置信息
		Interactive bool `flag:"desc=interactive mode;shorts=i"`
	}{}

	// AppExtAddCmd 添加应用扩展添加命令
	AppExtAddCmd = &gcli.Command{
		Name: "add",
		Desc: "Add a new kite extension",
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&extAddOpts)
			c.AddArg("name", "extension name", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			ext := kiteext.NewExt(c.Arg("name").String(), extAddOpts.Desc, extAddOpts.Path)

			// TODO set more info

			return app.Exts.Add(ext)
		},
	}
)
