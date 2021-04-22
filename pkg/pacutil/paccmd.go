package pacutil

import (
	"fmt"

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

type PacOpts struct {
	addr string
	file string
	gwUrl string
	gwfile string
	maxAge string
}

var pacOpts = PacOpts{}
var (
	PacServe = &gcli.Command{
		Name: "serve",
		Desc: "start an pac serve",
		Func: func(c *gcli.Command, args []string) error {
			return startServer(pacOpts)
		},

		Config: func(c *gcli.Command) {
			c.StrOpt(&pacOpts.addr, "addr", "a", ":11080", "server address")
			c.StrOpt(&pacOpts.file, "file", "f", "", "pac file path")
			c.FlagMeta("file").Required = true

			c.StrOpt(&pacOpts.gwfile, "gwfile", "", "", "gfw list file")
			c.StrOpt(&pacOpts.maxAge, "max-age", "m", "31536000", "Cache Control max-age")
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

	// example: pacgo catgw -f tmp/gfwlist-210422.txt
	GFWListCat = &gcli.Command{
		Name: "catgw",
		Desc: "decode gfw list content and print it",
		Func: func(c *gcli.Command, args []string) error {
			dst, err := DecodeGfwList(pacOpts.gwfile)
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
		Examples: `
{$fullCmd} -f ./tmp/gfwlist-210422.txt
`,
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
