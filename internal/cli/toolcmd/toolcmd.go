package toolcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/cli/appcmd"
	"github.com/inhere/kite-go/internal/cli/fscmd"
	"github.com/inhere/kite-go/internal/cli/pkgcmd"
	"github.com/inhere/kite-go/internal/cli/syscmd"
	"github.com/inhere/kite-go/internal/cli/taskcmd"
	"github.com/inhere/kite-go/internal/cli/textcmd"
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
		DateCmd,
		QuickJumpCmd,
		convcmd.ConvBaseCmd,
		// RunAnyCmd,
		// ScriptCmd,
		// RandomCmd,
		MathCalcCmd,
		// ScriptCmd,
		TranslateCmd,
		convcmd.NewTime2dateCmd(),
		convcmd.NewDate2tsCmd(),
		convcmd.NewConvPathSepCmd(),
		syscmd.NewClipboardCmd(),
		pkgcmd.PkgManageCmd,
		doccmd.DocumentCmd,
		doccmd.NewCheatCmd(),
		mdcmd.MkDownCmd,
		syscmd.NewQuickOpenCmd(),
		taskcmd.TaskManageCmd,
		// jsoncmd.JSONToolCmd,
		textcmd.NewMd5Cmd(),
		textcmd.NewHashCmd(),
		textcmd.NewUuidCmd(),
	},
	Config: func(c *gcli.Command) {

	},
}

// TranslateCmd command
var TranslateCmd = &gcli.Command{
	Name:    "translate",
	Aliases: []string{"trans", "tr", "fy"},
	Desc:    "translate the input contents to the target language",
	Config: func(c *gcli.Command) {
		// random string(number,alpha,), int(range)
	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
