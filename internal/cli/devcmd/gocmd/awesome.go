package gocmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
	"github.com/gookit/greq"
	"github.com/inhere/kite/internal/app"
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
		lang    string
		starMin int
		starMax int
		limit   int
		update  bool
		cnDocs  bool
		active  bool
		// keywords list
		keywords []string
	}{}

	langs = []string{"en", "zh-CN"}
)

// AwesomeGoCmd command
var AwesomeGoCmd = &gcli.Command{
	Name:    "awesome",
	Desc:    "view or search package on awesome go contents",
	Help:    "contents from: " + awesomeCnUrl,
	Aliases: []string{"awe"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&agOpts.cnDocs, "cn-doc", "cn", false, "the package contains Chinese readme")
		c.BoolOpt(&agOpts.update, "update", "up,u", false, "update the cached contents to latest")
		c.BoolOpt(&agOpts.active, "active", "", false, "the package status should be Active[update in the last week]")

		c.StrOpt2(&agOpts.lang, "lang,l", "language for the awesome-go contents, allow: en,zh-CN", func(opt *gflag.CliOpt) {
			opt.DefVal = "zh-CN"
			opt.Validator = func(val string) error {
				return errorx.IsIn(val, langs, "lang value must be in %v", langs)
			}
		})
		c.IntOpt2(&agOpts.starMin, "star-min,min", "limit the min star number")
		c.IntOpt2(&agOpts.starMax, "star-max,max", "limit the max star number")
		c.IntOpt(&agOpts.limit, "limit", "size", 50, "limit the max package number")

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

			if agOpts.update {
				c.Infoln("update completed")
				return nil
			}
		}

		keywords := c.Arg("keywords").Strings()

		agOpts.keywords = append(keywords, "github.com")
		c.Infoln("Search package by keywords:", agOpts.keywords)
		c.Println("----------------------- RESULT ----------------------")

		var sb strings.Builder
		sb.Grow(512)

		var category string

		var n int
		s := fsutil.LineScanner(cacheFile)
		for s.Scan() {
			// * [gookit/color](https://github.com/gookit/color) **star:782** Terminal color rendering tool library, support 16 colors, 256 colors, RGB color rendering output,
			// compatible with Windows.   [![There was an update last week][Green]](https://github.com/gookit/color)
			// line := s.Text()
			trimmed := strings.TrimSpace(s.Text())
			if strings.HasPrefix(trimmed, "##") {
				category = trimmed
			}

			// if !strings.HasPrefix(trimmed, "* ") {
			if !strings.HasPrefix(trimmed, "- ") {
				continue
			}

			if idx := strings.Index(trimmed, "[!["); idx > 0 {
				// pkgInfo := parsePkgInfo(trimmed[idx:]) // TODO
				trimmed = trimmed[:idx]
			}

			// match line
			if !textutil.IsMatchAll(trimmed, agOpts.keywords) {
				continue
			}

			// star number
			starNum := getStarNum(trimmed)
			if agOpts.starMin > 0 && starNum < agOpts.starMin {
				continue
			}
			if agOpts.starMax > 0 && starNum > agOpts.starMax {
				continue
			}

			if category != "" {
				c.Warnln(category)
				category = ""
			}

			n++
			fmt.Println(trimmed)

			if agOpts.limit > 0 && n >= agOpts.limit {
				break
			}
		}

		fmt.Println()
		c.Infoln("> Matched package result size:", n)
		// fmt.Println(sb.String())
		return nil
	},
}

func matchAwesomePkg() {
	// TODO
}

// AwesomePkgInfo struct
type AwesomePkgInfo struct {
	// GoDoc string // godoc document links
	CnDoc      bool // Contains Chinese documents
	WeekActive bool // There was an update last week
	YearActive bool // It hasn't been updated in the last year
	Archived   bool // The project has been archived
}

// parse Pkg Info
// [Awesome]: star > 2000
// [Green]: There was an update last week
// [Yellow]: It hasn't been updated in the last year
// [CN]: Contains Chinese documents
// [Archived]: The project has been archived
// [GoDoc]: godoc document links
func parsePkgInfo(str string) *AwesomePkgInfo {
	pi := &AwesomePkgInfo{}

	if strings.Contains(str, "[CN]") {
		pi.WeekActive = true
	}

	if strings.Contains(str, "[Green]") {
		pi.WeekActive = true
	} else if strings.Contains(str, "[Yellow]") {
		pi.YearActive = true
	} else if strings.Contains(str, "[Archived]") {
		pi.Archived = true
	}

	return pi
}

var matchStar = regexp.MustCompile(`\*star:(\d+)\*`)

func getStarNum(line string) int {
	// match like **star:778**
	ss := matchStar.FindStringSubmatch(line)
	if len(ss) > 1 {
		return strutil.Int2(ss[1])
	}

	return 0
}
