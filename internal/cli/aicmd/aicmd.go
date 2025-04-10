package aicmd

import (
	"github.com/gookit/gcli/v3"
)

var AICommand = &gcli.Command{
	Name: "ai",
	Desc: "AI tool command",
	Subs: []*gcli.Command{
		NewTranslateCmd(),
	},
}
