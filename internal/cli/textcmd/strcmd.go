package textcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/toolcmd/convcmd"
)

// TextToolCmd instance
var TextToolCmd = &gcli.Command{
	Name:    "text",
	Desc:    "useful commands for handle string text",
	Aliases: []string{"txt", "str", "string"},
	Subs: []*gcli.Command{
		StrCountCmd,
		StrSplitCmd,
		NewStrMatchCmd(),
		NewTextParseCmd(),
		NewTextSearchCmd(),
		NewReplaceCmd(),
		NewMd5Cmd(),
		NewHashCmd(),
		NewUuidCmd(),
		NewRandomStrCmd(),
		NewStringJoinCmd(),
		NewTemplateCmd(false),
		convcmd.NewTime2dateCmd(),
	},
}
