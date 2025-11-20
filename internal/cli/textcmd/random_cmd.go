package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
)

// NewRandomStrCmd create command
func NewRandomStrCmd() *gcli.Command {
	var randStrOpt = struct {
		number  int
		randTyp string
	}{}

	return &gcli.Command{
		Name:    "random",
		Aliases: []string{"rand"},
		Desc:    "quick generate random string or number",
		Config: func(c *gcli.Command) {
			// random string(number,alpha,), int(range)
			c.StrOpt(&randStrOpt.randTyp, "type", "t", "an", `The type of random string.
allow: num/number, a/alpha, hex, b64/base64, uuid, an/alpha_num`)
			c.IntOpt(&randStrOpt.number, "number", "n", 1, "The number of random string")

			c.AddArg("length", "The length of random string. allow: 1-1024").WithDefault("16")
		},
		Func: func(c *gcli.Command, _ []string) error {
			var str string
			ln := c.Arg("length").Int()

			for i := 0; i < randStrOpt.number; i++ {
				switch randStrOpt.randTyp {
				case "num", "number":
					fmt.Println("TODO")
				case "a", "alpha":
					str = strutil.RandWithTpl(ln, strutil.AlphaBet)
				case "hex":
					str = strutil.RandWithTpl(ln, strutil.HexChars)
				case "b64", "base64":
					str = strutil.RandWithTpl(ln, strutil.Base64Chars)
				case "an", "alpha_num":
					str = strutil.RandWithTpl(ln, strutil.AlphaNum2)
				default:
					return errorx.Ef("invalid type: %s", randStrOpt.randTyp)
				}
				fmt.Println(str)
			}
			return nil
		},
	}
}
