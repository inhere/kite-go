package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite/internal/cli/appcmd"
	"github.com/inhere/kite/internal/cli/devcmd/jsoncmd"
	"github.com/inhere/kite/internal/cli/fscmd"
	"github.com/inhere/kite/internal/cli/pkgcmd"
	"github.com/inhere/kite/internal/cli/syscmd"
	"github.com/inhere/kite/internal/cli/toolcmd/doctool"
	"github.com/inhere/kite/internal/cli/toolcmd/mdcmd"
	"github.com/inhere/kite/internal/cli/toolcmd/swagcmd"
)

// ToolsCmd command
var ToolsCmd = &gcli.Command{
	Name:    "tool",
	Aliases: []string{"tools"},
	Desc:    "provide some useful help tools commands",
	Subs: []*gcli.Command{
		swagcmd.SwaggerCmd,
		// textcmd.TextOperateCmd,
		syscmd.NewBatchRunCmd(),
		syscmd.NewEnvInfoCmd(),
		appcmd.NewPathMapCmd(),
		fscmd.NewFileCatCmd(),
		AutoJumpCmd,
		ConvBaseCmd,
		// RunAnyCmd,
		// ScriptCmd,
		RandomCmd,
		MathCalcCmd,
		Time2dateCmd,
		syscmd.NewClipboardCmd(),
		pkgcmd.PkgManageCmd,
		doctool.DocumentCmd,
		mdcmd.MkDownCmd,
		syscmd.NewQuickOpenCmd(),
		jsoncmd.JSONToolCmd,
	},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
