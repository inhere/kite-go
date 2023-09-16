package textcmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
)

// NewRandomStrCmd create command
func NewRandomStrCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "random",
		Aliases: []string{"rand"},
		Desc:    "quick generate random string or number",
		Config: func(c *gcli.Command) {
			// random string(number,alpha,), int(range)
			c.AddArg("length", "The length of random string. allow: 1-1024")
			c.AddArg("type", `The type of random string.
allow: num/number, a/alpha, hex, b64/base64, uuid, an/alpha_num`)
		},
		Func: func(c *gcli.Command, _ []string) error {
			n := c.Arg("length").Int()

			var str string
			switch c.Arg("type").String() {
			case "num", "number":
				fmt.Println("TODO")
			case "a", "alpha":
				str = strutil.RandWithTpl(n, strutil.AlphaBet)
			case "hex":
				str = strutil.RandWithTpl(n, strutil.HexChars)
			case "b64", "base64":
				str = strutil.RandWithTpl(n, strutil.Base64Chars)
			case "an", "alpha_num":
				str = strutil.RandWithTpl(n, strutil.AlphaNum2)
			default:
				return errorx.E("invalid type")
			}

			fmt.Println(str)
			return nil
		},
	}
}
