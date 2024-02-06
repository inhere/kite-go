package extcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/envutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// PlugCmd for kite plugins
var PlugCmd = &gcli.Command{
	Name:   "plug",
	Hidden: true,
	Desc:   "manage and run kite plugins(powered by traefik/yaegi interpreter)",
	Subs:   []*gcli.Command{},
	Config: func(c *gcli.Command) {
		c.AddArg("name", "input plugin name or path for execute")
	},
	Func: func(c *gcli.Command, args []string) error {
		name := c.Arg("name").String()
		if len(name) == 0 {
			return c.ShowHelp()
		}

		// dump.P(envutil.Getenv("GOPATH") + "/pkg/mod")

		// create new interpreter
		i := interp.New(interp.Options{
			GoPath: envutil.Getenv("GOPATH"),
			// Args: []string{},
		})
		if err := i.Use(stdlib.Symbols); err != nil {
			return err
		}

		// find plugin file
		fPath := apputil.ResolvePath(name)

		// run plugin file
		v, err := i.EvalPath(fPath)
		if err != nil {
			return err
		}

		dump.P(v.Interface())
		return nil
	},
}
