package pacutil

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/gookit/gcli/v3"
)

// refer links
// https://github.com/100apps/ipac
// https://zh.wikipedia.org/zh/%E4%BB%A3%E7%90%86%E8%87%AA%E5%8A%A8%E9%85%8D%E7%BD%AE
var PacTools = &gcli.Command{
	Name: "pac",
	// Hidden: true,
	Desc: "pac tools",
	Subs: []*gcli.Command{
		PacServe,
		GFWListCat,
		GFWList2pac,
		GFWListUpdate,
	},
}

var pacOpts = struct {
	addr string
	file string
	mAge string
	gwUrl string
	gwfile string
}{}

var (
	PacServe = &gcli.Command{
		Name: "serve",
		Desc: "start an pac serve",
		Func: func(c *gcli.Command, args []string) error {
			return startServer(pacOpts.addr, pacOpts.file, pacOpts.mAge)
		},

		Config: func(c *gcli.Command) {
			c.StrOpt(&pacOpts.addr, "addr", "a", ":11080", "server address")
			c.StrOpt(&pacOpts.file, "file", "f", "", "pac file path")
			c.FlagMeta("file").Required = true

			c.StrOpt(&pacOpts.mAge, "max-age", "m", "31536000", "Cache Control max-age")
		},
		Examples: `
{$fullCmd} -f ./tmp/gfwlist-210422.pac
`,
	}

	GFWListUpdate = &gcli.Command{
		Name: "gwup",
		Desc: "start an pac serve",
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},

		Aliases: []string{"upgw"},
	}

	// example: pac catgw tmp/gfwlist-210422.txt
	GFWListCat = &gcli.Command{
		Name: "catgw",
		Desc: "decode gfw list content and print it",
		Func: func(c *gcli.Command, args []string) error {

			src, err := ioutil.ReadFile(pacOpts.gwfile)
			if err != nil {
				return err
			}

			dst := make([]byte, len(src))
			_, err = base64.StdEncoding.Decode(dst, src)
			if err != nil {
				return err
			}

			fmt.Println(string(dst))
			return nil
		},

		Aliases: []string{"catgw"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&pacOpts.gwfile, "file", "f", "", "gfw list file")
			c.StrOpt(&pacOpts.gwUrl, "url", "u", "", "gfw list file url")
		},
	}

	GFWList2pac = &gcli.Command{
		Name: "gwconv",
		Desc: "convert gfw list to an pac file",
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},

		Aliases: []string{"convgw"},
	}
)
