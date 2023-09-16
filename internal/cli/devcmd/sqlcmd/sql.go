package sqlcmd

import "github.com/gookit/gcli/v3"

// SQLToolCmd instance
var SQLToolCmd = &gcli.Command{
	Name: "sql",
	Desc: "SQL tools",
	Subs: []*gcli.Command{
		Create2Mkdown,
		Conv2StructCmd,
		Insert2JSONCmd,
		NewCreate2JSONCmd(),
	},
}
