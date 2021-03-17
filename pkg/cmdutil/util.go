package cmdutil

import (
	"os/exec"

	"github.com/gookit/goutil/sysutil"
)

// OpenBrowser URL on browser
func OpenBrowser(url string)  {
	binName := "x-www-browser"
	if sysutil.IsWin() {
		binName = "start"
	} else if sysutil.IsMac() {
		binName = "open"
	}

	c := exec.Command(binName, url)
	if err := c.Run(); err != nil {
		panic(err)
	}
}
