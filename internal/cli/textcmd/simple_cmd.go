package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
)

// StrCountCmd instance
var StrCountCmd = &gcli.Command{
	Name:    "length",
	Aliases: []string{"len", "count"},
	Desc:    "count input string length, with rune length, utf8 length, text width",
	Config: func(c *gcli.Command) {
		c.AddArg("text", "input text contents for process. allow @c,@FILE", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(c.Arg("text").String())
		if err != nil {
			return err
		}

		fmt.Printf(
			"raw length: %d\n - rune len: %d\n - utf8 len: %d\n - width: %d\n",
			len(src),
			len([]rune(src)),
			strutil.Utf8Len(src),
			strutil.TextWidth(src),
		)
		return nil
	},
}
