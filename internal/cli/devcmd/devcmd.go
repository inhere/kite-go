package devcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/devcmd/gencmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/gocmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/javacmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/netcmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/phpcmd"
	"github.com/inhere/kite-go/internal/cli/devcmd/sqlcmd"
	jsoncmd "github.com/inhere/kite-go/internal/cli/devcmd/yamlcmd"
)

// DevToolsCmd command
var DevToolsCmd = &gcli.Command{
	Name:    "dev",
	Aliases: []string{"dt", "devtool"},
	Desc: "provide some useful develop tools commands",
	Subs: []*gcli.Command{
		HotReloadServe,
		IDEAToolCmd,
		ProjectCmd,
		gencmd.CodeGenCmd,
		javacmd.JavaToolCmd,
		gocmd.GoToolsCmd,
		phpcmd.PhpToolsCmd,
		sqlcmd.SQLToolCmd,
		jsoncmd.YamlToolCmd,
		netcmd.NewNMapCmd(),
		netcmd.NewNetcatCmd(),
		netcmd.NewTelnetClientCmd(),
		netcmd.NewTelnetServerCmd(),
	},
}
