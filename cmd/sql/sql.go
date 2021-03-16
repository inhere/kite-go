package sql

import "github.com/gookit/gcli/v3"

var SQLCmd = &gcli.Command{
	Name: "sql",
	Desc: "SQL tools",
	Subs: []*gcli.Command{
		Conv2Mkdown,
		Conv2Struct,
	},
}
