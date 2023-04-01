package sqlcmd

import "github.com/gookit/gcli/v3"

// SQLToolCmd instance
var SQLToolCmd = &gcli.Command{
	Name: "sql",
	Desc: "SQL tools",
	Subs: []*gcli.Command{
		Conv2Mkdown,
		Conv2StructCmd,
		Conv2JSONCmd,
	},
}
