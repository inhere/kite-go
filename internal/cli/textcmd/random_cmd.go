package textcmd

import (
	"encoding/base64"
	"fmt"

	"github.com/gofrs/uuid/v5"
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
				str = base64.StdEncoding.EncodeToString([]byte(strutil.RandWithTpl(n, strutil.Base64Chars)))
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

// NewUuidCmd create command
func NewUuidCmd() *gcli.Command {
	var uuidOpts = struct {
		V1 bool `flag:"desc=Generate UUIDv1 (timestamp based);shorts=1"`
		// V2  bool `flag:"desc=Generate UUIDv2 (DCE Security version, based on POSIX UID/GID);shorts=2"` // not support
		V3  bool `flag:"desc=Generate UUIDv3 (name-based, MD5);shorts=3"`
		V4  bool `flag:"desc=Generate UUIDv4 (random);shorts=4"`
		V5  bool `flag:"desc=Generate UUIDv5 (name-based, SHA-1);shorts=5"`
		V6  bool `flag:"desc=[draft]Generate UUIDv6 (ordered-time);shorts=6"`
		V7  bool `flag:"desc=[draft]Generate UUIDv7 (ordered-time, DCE Security version, based on POSIX UID/GID);shorts=7"`
		ver uint8
		Ver uint `flag:"desc=Want generate UUID version, allow 1-7;shorts=v;default=4"`
		Num int  `flag:"desc=Generate the given number of UUIDs;shorts=n;default=1"`
		// args for v3, v5
		ns   uuid.UUID
		Ns   string `flag:"desc=The namespace string for v3, v5"`
		Name string `flag:"desc=The name string for v3, v5"`
	}{}

	return &gcli.Command{
		Name:    "uuid",
		Aliases: []string{"uid"},
		Desc:    "quick generate a UUID string",
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&uuidOpts)
		},
		Func: func(c *gcli.Command, _ []string) error {
			if uuidOpts.Num <= 0 {
				uuidOpts.Num = 1
			}

			if uuidOpts.V1 {
				uuidOpts.ver = uuid.V1
			} else if uuidOpts.V3 {
				uuidOpts.ver = uuid.V3
			} else if uuidOpts.V4 {
				uuidOpts.ver = uuid.V4
			} else if uuidOpts.V5 {
				uuidOpts.ver = uuid.V5
			} else if uuidOpts.V6 {
				uuidOpts.ver = uuid.V6
			} else if uuidOpts.V7 {
				uuidOpts.ver = uuid.V7
			} else {
				uuidOpts.ver = uint8(uuidOpts.Ver)
			}

			if uuidOpts.ver == uuid.V3 || uuidOpts.ver == uuid.V5 {
				uuidOpts.Name = strutil.OrElse(uuidOpts.Name, "github.com/inhere/kite-go")
				if uuidOpts.Ns == "" {
					uuidOpts.ns = uuid.NamespaceDNS
				} else {
					uuidOpts.ns = uuid.Must(uuid.FromString(uuidOpts.Ns))
				}
			}

			var uid uuid.UUID
			for i := 0; i < uuidOpts.Num; i++ {
				switch uuidOpts.ver {
				case uuid.V1:
					uid = uuid.Must(uuid.NewV1())
				case uuid.V3:
					uid = uuid.NewV3(uuidOpts.ns, uuidOpts.Name)
				case uuid.V4:
					uid = uuid.Must(uuid.NewV4())
				case uuid.V5:
					uid = uuid.NewV5(uuidOpts.ns, uuidOpts.Name)
				case uuid.V6:
					uid = uuid.Must(uuid.NewV6())
				case uuid.V7:
					uid = uuid.Must(uuid.NewV7())
				default:
					return errorx.Errf("invalid UUID version %d", uuidOpts.Ver)
				}

				fmt.Println(uid.String())
			}
			return nil
		},
	}
}
