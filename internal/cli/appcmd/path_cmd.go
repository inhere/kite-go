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
		c.AddArg("name", "special path name on the kite, allow: base, config, tmp")
	},
	Func: func(c *gcli.Command, args []string) error {
		if kpOpts.list {
			show.AList("Kite paths", app.App().Config)
			return nil
		}

		name := c.Arg("name").String()
		if name == "" {
			return errorx.Raw("please input name for show path")
		}

		var path = app.App().PathByName(name)
		if path == "" {
			return errorx.Rawf("not found path for %q", name)
		}

		fmt.Println(path)
		return nil
	},
}

// paCmdOpts struct
type paCmdOpts struct {
	List bool `flag:"list all path alias map;;;l"`
	Name string
}

// NewPathMapCmd command
func NewPathMapCmd() *gcli.Command {
	var paOpts = &paCmdOpts{}

	return &gcli.Command{
		Name:    "pathmap",
		Aliases: []string{"path-alias", "userpath", "pmap"},
		Desc:    "show user custom path mapping in app(config:path_map)",
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.UseSimpleRule().FromStruct(paOpts))
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
