package swagger

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/sysutil"
)

var docGenOpts = struct {
	Output string
}{}

var DocGen = &gcli.Command{
	Name:   "gen",
	Desc: "generate swagger doc files by package: swaggo/swag",
	Config: func(c *gcli.Command) {
		c.StrOpt(&docGenOpts.Output, "output", "o", "./static", `the output directory for generated doc files`)
	},
	Func: func(c *gcli.Command, _ []string) error {
		// swag init -o static
		// rm static/docs.go
		ret, err := sysutil.ExecCmd("swag", []string{
			"init",
			"-o",
			docGenOpts.Output,
		})
		fmt.Println(ret)

		return err
	},
	Help: `
Install swag:
    go get -u github.com/swaggo/swag/cmd/swag
`,
}
