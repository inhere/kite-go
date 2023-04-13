package doccmd

import (
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/stdio"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/greq"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
)

const RandTopic = ":random"
const ChtHost = "https://cht.sh/"
const chtHelp = `
<b>Special pages</b>
There are several special pages that are not cheat sheets. Their names start with colon and have special meaning.

Getting started:

    :help               description of all special pages and options
    :intro              cheat.sh introduction, covering the most important usage questions
    :list               list all cheat sheets (can be used in a subsection too: /go/:list)

Command line client cht.sh and shells support:

    :cht.sh             code of the cht.sh client
    :bash_completion    bash function for tab completion
    :bash               bash function and tab completion setup
    :fish               fish function and tab completion setup
    :zsh                zsh function and tab completion setup

Editors support:

    :vim                cheat.sh support for Vim
    :emacs              cheat.sh function for Emacs
    :emacs-ivy          cheat.sh function for Emacs (uses ivy)

Other pages:

    :post               how to post new cheat sheet
    :styles             list of color styles
    :styles-demo        show color styles usage examples
    :random             fetches a random page (can be used in a subsection too: /go/:random)
`

var chtOpt = struct {
	Refresh bool   `flag:"ignore cached result, re-request remote cheat server;;;r"`
	Search  string `flag:"search the cheat sheets by keywords;;;s"`
	// List cached cheat sheets
	List bool `flag:"list all cached cheat sheets in local;;;l"`
	// Topic    string
	// Question []string
}{}

// NewCheatCmd instance
//
// Query cheat for development
//
//   - github: https://github.com/chubin/cheat.sh
//
// Example:
//
//	curl cheat.sh/tar
//	curl cht.sh/curl
//	curl https://cheat.sh/rsync
//	curl https://cht.sh/php
//
//	curl cht.sh/go/:list
//	curl cht.sh/go/reverse+a+list
//	curl cht.sh/python/random+list+elements
//	curl cht.sh/js/parse+json
//	curl cht.sh/lua/merge+tables
//	curl cht.sh/clojure/variadic+function
func NewCheatCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "cheat",
		Aliases: []string{"cht", "chtsh"},
		Desc:    "Query cheat for development. from https://cht.sh",
		Help:    chtHelp,
		Examples: `
{$fullCmd} php strlen
{$fullCmd} go reverse list
{$fullCmd} java Optional
`,
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&chtOpt, gflag.TagRuleSimple)
			c.AddArg("topic", "the topic name. e.g: php, go, java, python, js, lua, clojure")
			c.AddArg("question", "The questions on the topic. e.g: strlen, reverse list, Optional", false, true)
		},
		Func: func(c *gcli.Command, _ []string) error {
			if chtOpt.Search != "" {
				searchApi := ChtHost + "~" + chtOpt.Search
				colorp.Infoln("will send request to", searchApi)

				reqOpt := greq.NewOpt(func(opt *greq.Option) {
					opt.SetUserAgent(greq.AgentCURL)
				})
				resp, err := greq.GetDo(searchApi, reqOpt)
				if err != nil {
					return err
				}

				colorp.Infoln("RESULT:")
				stdio.Writeln(resp.BodyString())
			}

			topic := c.Arg("topic").String()
			if topic == "" {
				return errorx.Rawf("missing topic name for query. e.g: php, go, java, python, js, lua, clojure")
			}

			queries := c.Arg("question").Strings()
			result, err := queryCheat(topic, queries, chtOpt.Refresh)
			if err != nil {
				return err
			}

			colorp.Infoln("RESULT:")
			stdio.Writeln(result)
			return nil
		},
	}
}

func queryCheat(topic string, queries []string, refresh bool) (string, error) {
	var cacheDir string
	if setDir := app.Cfg().String("cheat.cache_dir"); setDir != "" {
		cacheDir = apputil.ResolvePath(setDir)
	} else {
		cacheDir = app.App().TmpPath("cheat")
	}

	queryStr := strings.Join(queries, "/")
	cacheFile := fsutil.JoinSubPaths(cacheDir, topic)
	if topic[0] == ':' {
		cacheFile += ".txt"
	} else if queryStr != "" {
		cacheFile += "/" + queryStr + ".txt"
	} else {
		cacheFile += "/_topic.txt"
	}

	if !refresh && fsutil.IsFile(cacheFile) {
		colorp.Infoln("use cached result:", cacheFile)
		return fsutil.ReadStringOrErr(cacheFile)
	}

	chtApiUrl := ChtHost + topic
	if queryStr != "" {
		chtApiUrl += "/" + queryStr
	}

	colorp.Infoln("will send request to", chtApiUrl)
	reqOpt := greq.NewOpt(func(opt *greq.Option) {
		opt.SetUserAgent(greq.AgentCURL)
	})
	resp, err := greq.GetDo(chtApiUrl, reqOpt)
	if err != nil {
		return "", err
	}

	bodyStr := resp.BodyString()

	// not found
	if resp.ContentLength < 300 &&
		(strings.Contains(bodyStr, "Unknown topic.") || strings.Contains(bodyStr, "Unknown cheat sheet")) {
		return bodyStr, nil
	}

	// an random document.
	if topic == RandTopic {
		fline := strutil.FirstLine(bodyStr)
		name := strings.Trim(color.ClearCode(fline), "#/ \t\n\r\x0B")
		name = strings.TrimPrefix(name, "cheat:")
		if len(name) == 0 {
			return bodyStr, nil
		}

		colorp.Infoln("found the random document:", name)
		cacheFile = cacheDir + "/random/" + name + ".txt"
	}

	if resp.ContentLength > 0 {
		colorp.Infoln("write cache file:", cacheFile)
		_, err := fsutil.PutContents(cacheFile, bodyStr)
		if err != nil {
			colorp.Errorln("write cache file failed:", err)
		}
	}

	return bodyStr, nil
}
