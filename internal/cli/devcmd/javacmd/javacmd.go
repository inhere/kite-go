package javacmd

import (
	"github.com/gookit/gcli/v3"
)

// JavaToolCmd command
var JavaToolCmd = &gcli.Command{
	Name: "java",
	// Aliases: []string{},
	Desc: "provide some useful java tools commands",
	Subs: []*gcli.Command{},
	Config: func(c *gcli.Command) {

	},
	// Func: func(c *gcli.Command, _ []string) error {
	// 	return errors.New("TODO")
	// },
}
