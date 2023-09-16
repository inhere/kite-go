package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/app"
)

var kpOpts = struct {
	list bool
}{}

// KitePathCmd command
var KitePathCmd = &gcli.Command{
	Name:    "path",
	Aliases: []string{"paths"},
	Desc:    "show the kite system path information",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&kpOpts.list, "list, all, a, l", "display all paths for the kite")
		c.AddArg("path", "special alias path for resolve. prefix: base, config, tmp")
	},
	Func: func(c *gcli.Command, args []string) error {
		if kpOpts.list {
			show.AList("Kite paths", app.App().Config)
			return nil
		}

		path := c.Arg("path").String()
		if path == "" {
			return errorx.Raw("please input path for resolve")
		}

		path = app.App().ResolvePath(path)
		if path == "" {
			return errorx.Rawf("not found path for %q", path)
		}

		fmt.Println(path)
		return nil
	},
}

// paCmdOpts struct
type paCmdOpts struct {
	List bool `flag:"desc=list all user path alias map;shorts=a,l"`
	Name string
}

// NewPathMapCmd command
func NewPathMapCmd() *gcli.Command {
	var paOpts = &paCmdOpts{}

	return &gcli.Command{
		Name:    "pathmap",
		Aliases: []string{"path-alias", "pmap", "pmp"},
		Desc:    "show user custom path mapping in kite(config:path_map)",
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.FromStruct(paOpts))
			c.AddArg("name", "get path of the input alias name").WithAfterFn(func(a *gflag.CliArg) error {
				paOpts.Name = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			if paOpts.List {
				show.AList("User Paths:", app.PathMap.Data())
				return nil
			}

			if paOpts.Name != "" {
				fmt.Println(app.PathMap.Resolve(paOpts.Name))
				return nil
			}
			return errorx.New("please input name for get path")
		},
	}
}
