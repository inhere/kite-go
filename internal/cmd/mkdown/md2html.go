package mkdown

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	gmhtml "github.com/gomarkdown/markdown/html"
	gmparser "github.com/gomarkdown/markdown/parser"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gdparser "github.com/yuin/goldmark/parser"
	gdhtml "github.com/yuin/goldmark/renderer/html"
)

// filetypes: [".md", ".markdown", ".mdown"]
type md2html struct {
	cmd *gcli.Command
	// options
	toc   bool
	page  bool
	latex bool

	tocOnly   bool
	fractions bool

	smartyPants bool
	latexDashes bool
	// Sets HTML output to a simple form:
	//  - No head
	//  - No body tags
	//  - ids, classes, and style are stripped out
	//  - Just bare minimum HTML tags and attributes
	//  - extension modifications included
	htmlSimple bool

	css string
	// driver:
	// gm    gomarkdown
	// bf 	 blackfriday
	driver string
	output string
	// "markdown", "github", "gitlab"
	style string
}

const (
	defaultTitle = ""
	// driver names
	driverBF = "bf"
	driverGM = "gm"
)

var mh = md2html{}

var (
	drivers = map[string]string{
		"bf": "blackfriday",
		"gm": "gomarkdown",
		"gd": "goldmark",
	}
)

/*
DOC: https://developer.github.com/v3/markdown/#render-an-arbitrary-markdown-document
curl https://api.github.com/markdown/raw -X "POST" -H "Content-Type: text/plain" -d "Hello world github/linguist#1 **cool**, and #1!"

DOC: https://docs.gitlab.com/ee/api/markdown.html#render-an-arbitrary-markdown-document
curl --header Content-Type:application/json --data '{"text":"Hello world! :tada:", "gfm":true, "project":"group_example/project_example"}' https://gitlab.example.com/api/v4/markdown
*/

// Markdown2HTML Convert Markdown to HTML
// styles from https://github.com/facelessuser/MarkdownPreview
//
// "image_path": "https://github.githubassets.com/images/icons/emoji/unicode/",
// "non_standard_image_path": "https://github.githubassets.com/images/icons/emoji/"
var Markdown2HTML = &gcli.Command{
	Name:    "html",
	Desc:  "convert one or multi markdown file to html",
	Aliases: []string{"tohtml"},
	// Config:  nil,
	// Examples: "",
	Func: mh.Handle,
	Config: func(c *gcli.Command) {
		c.BoolOpt(&mh.toc, "toc", "", false,
			"Generate a table of contents (implies --latex=false)")
		flag.BoolVar(&mh.tocOnly, "toconly", false,
			"Generate a table of contents only (implies -toc)")
		c.BoolOpt(&mh.page, "page", "", false,
			"Generate a standalone HTML page (implies --latex=false)")
		c.BoolOpt(&mh.latex, "latex", "", false,
			"Generate LaTeX output instead of HTML")

		c.BoolOpt(&mh.smartyPants, "smartypants", "", true,
			"Apply smartypants-style substitutions")
		c.BoolOpt(&mh.latexDashes, "latexdashes", "", true,
			"Use LaTeX-style dash rules for smartypants")
		c.BoolOpt(&mh.fractions, "fractions", "", true,
			"Use improved fraction rules for smartypants")
		c.BoolOpt(&mh.htmlSimple, "html-simple", "", true,
			"Sets HTML output to a simple, just bare minimum HTML tags and attributes")

		c.StrOpt(&mh.css, "css", "", "",
			"Link to a CSS stylesheet (implies --page)")
		c.StrOpt(&mh.output, "output", "", "",
			"the rendered content output, default output STDOUT")
		c.StrOpt(&mh.driver, "driver", "", "bf",
			`set the markdown renderer driver.
allow:
bf - blackfriday
gm - gomarkdown
gd - goldmark
`)

		c.AddArg("files", "the listed files will be render to html", false, true)

		// save
		mh.cmd = c
	},
}

func (mh md2html) Handle(c *gcli.Command, args []string) (err error) {
	// enforce implied options
	if mh.css != "" {
		mh.page = true
	}
	if mh.page {
		mh.latex = false
	}
	if mh.toc {
		mh.latex = false
	}

	color.Info.Println("Work Dir:", c.WorkDir())
	color.Info.Println("Use Driver:", mh.driverName())

	mdString := `
# title

## h2

hello

### h3
`

	if mh.driver == driverBF {
		err = mh.blackFriday([]byte(mdString), args)
	} else {
		err = mh.goMarkdown([]byte(mdString), args)
	}

	// color.Success.Println("Complete")
	return
}

func (mh md2html) driverName() string {
	if name, ok := drivers[mh.driver]; ok {
		return name
	}

	return drivers[driverBF]
}

func (mh md2html) blackFriday(input []byte, args []string) (err error) {
	// blackfriday.Run()

	// return mh.outToWriter(buf.Bytes())
	return errors.New("TODO: current not support 'blackFriday'")
}

func (mh md2html) goldMark(source []byte, args []string) (err error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			gdparser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gdhtml.WithHardWraps(),
			gdhtml.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return err
	}

	return mh.outToWriter(buf.Bytes())
}

func (mh md2html) goMarkdown(input []byte, args []string) (err error) {
	// set up options
	var extensions = gmparser.NoIntraEmphasis |
		gmparser.Tables |
		gmparser.FencedCode |
		gmparser.Autolink |
		gmparser.Strikethrough |
		gmparser.SpaceHeadings

	var renderer markdown.Renderer
	if mh.latex {
		// render the data into LaTeX
		// renderer = markdown.LatexRenderer(0)
		color.Comment.Println("unsupported")
		return
	} else {
		// render the data into HTML
		var htmlFlags gmhtml.Flags
		// if xhtml {
		// 	htmlFlags |= html.UseXHTML
		// }
		if mh.smartyPants {
			htmlFlags |= gmhtml.Smartypants
		}
		if mh.fractions {
			htmlFlags |= gmhtml.SmartypantsFractions
		}
		if mh.latexDashes {
			htmlFlags |= gmhtml.SmartypantsLatexDashes
		}

		title := ""
		if mh.page {
			htmlFlags |= gmhtml.CompletePage
			title = getTitle(input)
		}
		if mh.toc {
			htmlFlags |= gmhtml.TOC
		}

		params := gmhtml.RendererOptions{
			Flags: htmlFlags,
			Title: title,
			CSS:   mh.css,
		}
		renderer = gmhtml.NewRenderer(params)
	}

	// parse and render
	psr := gmparser.NewWithExtensions(extensions)

	htmlBts := markdown.ToHTML(input, psr, renderer)

	return mh.outToWriter(htmlBts)
}

func (mh md2html) outToWriter(htmlText []byte) (err error) {
	// output the result
	var out *os.File
	if mh.output == "" {
		color.Info.Println("OUTPUT:")
		out = os.Stdout
	} else {
		if out, err = os.Create(mh.output); err != nil {
			return fmt.Errorf("error creating %s: %v", mh.output, err)
		}
		defer out.Close()
	}

	if _, err = out.Write(htmlText); err != nil {
		err = fmt.Errorf("error writing output: %s", err.Error())
	}
	return
}

// try to guess the title from the input buffer
// just check if it starts with an <h1> element and use that
func getTitle(input []byte) string {
	i := 0

	// skip blank lines
	for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
		i++
	}
	if i >= len(input) {
		return defaultTitle
	}
	if input[i] == '\r' && i+1 < len(input) && input[i+1] == '\n' {
		i++
	}

	// find the first line
	start := i
	for i < len(input) && input[i] != '\n' && input[i] != '\r' {
		i++
	}
	line1 := input[start:i]
	if input[i] == '\r' && i+1 < len(input) && input[i+1] == '\n' {
		i++
	}
	i++

	// check for a prefix header
	if len(line1) >= 3 && line1[0] == '#' && (line1[1] == ' ' || line1[1] == '\t') {
		return strings.TrimSpace(string(line1[2:]))
	}

	// check for an underlined header
	if i >= len(input) || input[i] != '=' {
		return defaultTitle
	}
	for i < len(input) && input[i] == '=' {
		i++
	}
	for i < len(input) && (input[i] == ' ' || input[i] == '\t') {
		i++
	}
	if i >= len(input) || (input[i] != '\n' && input[i] != '\r') {
		return defaultTitle
	}

	return strings.TrimSpace(string(line1))
}
