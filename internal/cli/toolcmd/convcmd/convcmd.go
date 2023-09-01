package convcmd

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
)

var convBaseOpts = struct {
	From int `flag:"desc=from base type. allow: 2-64;shorts=f;default=10"`
	To   int `flag:"desc=to base type, if eq 1, will show: 2,10,16,32,36,62,64 base. allow: 1-64;shorts=t"`
}{}

// ConvBaseCmd command
// Base
// Binary
// decimal
// Base 8
var ConvBaseCmd = &gcli.Command{
	Name:    "conv-base",
	Aliases: []string{"base", "cb"},
	Desc:    "Convert base data type. eg: binary, decimal(10), base8, hex(16), base64, base32 ...",
	Config: func(c *gcli.Command) {
		// random string(number,alpha,), int(range)
		c.MustFromStruct(&convBaseOpts)
		c.AddArg("input", "want convert base string contents")
	},
	Func: func(c *gcli.Command, _ []string) (err error) {
		str := c.Arg("input").String()
		str, err = apputil.ReadSource(str)
		if err != nil {
			return err
		}

		c.Infoln("Input String:", str)
		c.Warnln("Conv Results:")
		commonBases := []int{2, 10, 16, 32, 36, 62, 64}
		if convBaseOpts.To < 2 {
			for _, toBase := range commonBases {
				if toBase == convBaseOpts.From {
					continue
				}

				dst := strutil.BaseConv(str, convBaseOpts.From, toBase)
				fmt.Printf("convert to %d: %s\n", toBase, dst)
			}
		} else {
			dst := strutil.BaseConv(str, convBaseOpts.From, convBaseOpts.To)
			fmt.Printf("convert to %d: %s\n", convBaseOpts.To, dst)
		}

		return
	},
}

var cpsOpts = struct {
	format string
}{}

// NewConvPathSepCmd instance
func NewConvPathSepCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "conv-path",
		Aliases: []string{"conv-sep"},
		Desc:    "Quick convert unix path to Windows path",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&cpsOpts.format, "format, f", `sets the target format, will auto-detect on is empty.
allow: w/win/windows, l/lin/linux/unix`)
			c.AddArg("input", "want convert path contents")
		},
		Func: func(c *gcli.Command, _ []string) (err error) {
			pathStr := c.Arg("input").String()
			pathStr, err = apputil.ReadSource(pathStr)
			if err != nil {
				return err
			}

			// will auto-detect on is empty
			if cpsOpts.format == "" {
				// win -> linux
				if strings.ContainsRune(pathStr, '\\') {
					cpsOpts.format = "linux"
					// linux -> win
				} else if strings.ContainsRune(pathStr, '/') {
					cpsOpts.format = "win"
				}
			}

			switch strings.ToLower(cpsOpts.format) {
			case "w", "win", "windows":
				pathStr = fsutil.SlashPath(pathStr)
			case "l", "lin", "linux":
				pathStr = fsutil.UnixPath(pathStr)
			}
			fmt.Println(pathStr)
			return nil
		},
	}
}
