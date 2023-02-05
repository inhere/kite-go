package syscmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/sysutil"
)

// SysCmd command
var SysCmd = &gcli.Command{
	Name:    "sys",
	Aliases: []string{"os", "system"},
	Desc:    "provide some useful system commands",
	Subs: []*gcli.Command{
		SearchExeCmd,
		WhichExeCmd,
	},
}

// WhichExeCmd command
var WhichExeCmd = &gcli.Command{
	Name:    "which",
	Aliases: []string{"whereis", "type"},
	Desc:    "show full path for the executable name",
	Config: func(c *gcli.Command) {
		c.AddArg("name", "executable name for match", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		name := c.Arg("keyword").String()
		file, err := sysutil.Executable(name)
		if err != nil {
			return err
		}

		fmt.Println(file)
		return nil
	},
}

// SearchExeCmd command
var SearchExeCmd = &gcli.Command{
	Name:    "find-bin",
	Aliases: []string{"find-exe", "search"},
	Desc:    "search executable file in system PATH",
	Config: func(c *gcli.Command) {
		c.AddArg("keyword", "keywords for search in PATH dirs", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		files := sysutil.SearchPath(c.Arg("keyword").String(), 10)

		show.AList("Matched files:", files)
		return nil
	},
}
