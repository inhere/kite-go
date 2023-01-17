package comtool

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/comtool/swagger"
)

// ToolsCmd command
var ToolsCmd = &gcli.Command{
	Name:    "tool",
	Aliases: []string{"tools"},
	Desc:    "provide some useful tools commands",
	Subs: []*gcli.Command{
		swagger.SwaggerCmd,
		HttpServe,
		BatchRun,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
