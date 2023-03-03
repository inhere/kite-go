package appcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite/internal/app"
)

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
		Aliases: []string{"path-alias", "pmap"},
		Desc:    "show custom path aliases mapping in app(config:pathmap)",
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.UseSimpleRule().FromStruct(paOpts))
			c.AddArg("name", "get path of the input alias name").WithAfterFn(func(a *gflag.CliArg) error {
				paOpts.Name = a.String()
				return nil
			})
		},
		Func: func(c *gcli.Command, args []string) error {
			if paOpts.List {
				show.AList("Path aliases", app.PathMap.Data())
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
