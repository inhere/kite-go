package gocmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/greq"
	"github.com/inhere/kite/app"
)

var (
	// origin repo url
	awesomeGoUrl = "https://github.com/avelino/awesome-go"
	// https://raw.githubusercontent.com/yinggaozhen/awesome-go-cn/master/README_EN.md
	// https://raw.githubusercontent.com/yinggaozhen/awesome-go-cn/master/README.md
	awesomeCnUrl = "https://github.com/yinggaozhen/awesome-go-cn"

	awesomeENTextUrl = "https://raw.githubusercontent.com/yinggaozhen/awesome-go-cn/master/README_EN.md"
	awesomeZHTextUrl = "https://raw.githubusercontent.com/yinggaozhen/awesome-go-cn/master/README.md"

	// cmd options
	agOpts = struct {
		lang   string
		update bool
	}{}
)

// AwesomeGoCmd command
var AwesomeGoCmd = &gcli.Command{
	Name:    "awesome",
	Desc:    "view or search package on awesome go contents",
	Help:    "contents from: " + awesomeCnUrl,
	Aliases: []string{"awe"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&agOpts.update, "update", "up,u", false, "update the cached contents to latest")
		c.StrOpt(&agOpts.lang, "lang", "l", "en", "language for the awesome-go contents, allow: en,zh-CN")
		c.AddArg("keywords", "the keyword for search awesome-go contents", false, true)
	},
	Func: func(c *gcli.Command, args []string) error {
		mkdownUrl := awesomeENTextUrl
		cacheFile := app.App().PathResolve("$tmp/awesome-go-cn.EN.md")
		if agOpts.lang == "zh-CN" {
			mkdownUrl = awesomeZHTextUrl
			cacheFile = app.App().PathResolve("$tmp/awesome-go-cn.zh-CN.md")
		}

		if agOpts.update || !fsutil.FileExists(cacheFile) {
			err := greq.MustDo("GET", mkdownUrl).SaveFile(cacheFile)
			if err != nil {
				return err
			}
		}

		s := fsutil.LineScanner(cacheFile)
		for s.Scan() {
			fmt.Println(s.Text())
		}

		return nil
	},
}
