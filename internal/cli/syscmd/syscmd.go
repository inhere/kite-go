package syscmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/sysutil"
	"github.com/inhere/kite-go/internal/app"
)

// SysCmd command
var SysCmd = &gcli.Command{
	Name:    "sys",
	Aliases: []string{"os", "system"},
	Desc:    "provide some useful system commands",
	Subs: []*gcli.Command{
		SearchExeCmd,
		WhichExeCmd,
		SysInfoCmd,
		NewQuickOpenCmd(),
		NewBatchRunCmd(),
		NewEnvInfoCmd(),
		NewClipboardCmd(),
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
		name := c.Arg("name").String()
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

		show.AList("Matched exe files:", files)
		return nil
	},
}

// SysInfoCmd command
var SysInfoCmd = &gcli.Command{
	Name: "info",
	// Aliases: []string{"i"},
	Desc: "display current operation system information",
	Config: func(c *gcli.Command) {
		// c.AddArg("keyword", "keywords for search in PATH dirs", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		info := map[string]any{
			"os platform": runtime.GOOS,
			"os arch":     runtime.GOARCH,
		}

		if sysutil.IsWindows() {
			if sysutil.IsMSys() {
				info["hosts file1"] = "/etc/hosts"
				info["hosts file"] = "/c/Windows/System32/drivers/etc/hosts"
			} else {
				info["hosts file"] = "C:\\Windows\\System32\\drivers\\etc\\hosts"
			}
		} else {
			info["hosts file"] = "/etc/hosts"
		}

		show.AList("System info:", info)
		return nil
	},
}

// NewQuickOpenCmd command
func NewQuickOpenCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "open",
		Aliases: []string{"open-exe"},
		Desc:    "open input file or dir or remote URL address",
		Config: func(c *gcli.Command) {
			c.AddArg("name", "bin name or URL address", true)
		},
		Func: func(c *gcli.Command, _ []string) error {
			name := c.Arg("name").String()

			var dstFile = name
			if strings.Contains(name, "/") {
				// special github url
				if strings.HasPrefix(name, gitw.GitHubHost) {
					dstFile = "https://" + name
					// } else if fsutil.PathExists(name) {
					// 	// nothing ...
					// } else if validate.IsURL(name) {
				}
			} else if app.OpenMap.HasAlias(name) {
				dstFile = app.OpenMap.ResolveAlias(name)
			}

			c.Infoln("Will Open the:", dstFile)
			return sysutil.Open(dstFile)
		},
	}
}
