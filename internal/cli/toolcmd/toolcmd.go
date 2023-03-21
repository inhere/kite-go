package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/appcmd"
	"github.com/inhere/kite-go/internal/cli/fscmd"
	"github.com/inhere/kite-go/internal/cli/pkgcmd"
	"github.com/inhere/kite-go/internal/cli/syscmd"
	"github.com/inhere/kite-go/internal/cli/taskcmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd/convcmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd/doccmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd/mdcmd"
	"github.com/inhere/kite-go/internal/cli/toolcmd/swagcmd"
)

// ToolsCmd command
var ToolsCmd = &gcli.Command{
	Name:    "tool",
	Aliases: []string{"tools"},
	Desc:    "provide some useful help tools commands",
	Subs: []*gcli.Command{
		swagcmd.SwaggerCmd,
		syscmd.NewBatchRunCmd(),
		syscmd.NewEnvInfoCmd(),
		appcmd.NewPathMapCmd(),
		fscmd.NewFileCatCmd(),
		AutoJumpCmd,
		convcmd.ConvBaseCmd,
		// RunAnyCmd,
		// ScriptCmd,
		RandomCmd,
		MathCalcCmd,
		convcmd.NewTime2dateCmd(),
		convcmd.NewConvPathSepCmd(),
		syscmd.NewClipboardCmd(),
		pkgcmd.PkgManageCmd,
		doccmd.DocumentCmd,
		mdcmd.MkDownCmd,
		syscmd.NewQuickOpenCmd(),
		taskcmd.TaskManageCmd,
		// jsoncmd.JSONToolCmd,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
