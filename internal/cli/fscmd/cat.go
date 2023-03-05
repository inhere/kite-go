package fscmd

import (
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/glamour"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/stdio"
	"github.com/inhere/kite/pkg/kiteext"
)

var fcOpts = struct {
	clip   bool
	stdin  bool
	style  string
	format string
}{
	style: "auto",
}

// NewFileCatCmd instance
func NewFileCatCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "cat",
		Aliases: []string{"see", "bat"},
		Desc:    "view file contents like `cat`, but with style",
		Config: func(c *gcli.Command) {
			c.BoolOpt2(&fcOpts.stdin, "stdin, in, i", "read and cat the contents from stdin")
			c.BoolOpt2(&fcOpts.clip, "clip, c", "read and cat the contents from clipboard")

			styleNames := styles.Names()
			sort.Strings(styleNames)
			var sb strings.Builder
			for i, name := range styleNames {
				sb.WriteString(name)
				sb.WriteString(", ")
				if i+1%8 == 0 {
					sb.WriteByte('\n')
				}
			}
			c.StrOpt2(&fcOpts.style, "style, s", "sets the render style, default is auto.\n allow: auto, "+sb.String())
			c.StrOpt2(&fcOpts.format, "format, f", "sets the content format, default auto parse by filename")

			c.AddArg("files", "want cat file, allow multi files", false, true)
		},
		Func: fileCat,
	}
}

func fileCat(c *gcli.Command, _ []string) error {
	// format := strutil.OrElse(fcOpts.format, "markdown")
	format := fcOpts.format
	if fcOpts.stdin {
		return renderOneFile(kiteext.DstStdin, format)
	}

	if fcOpts.clip {
		return renderOneFile(kiteext.DstClip, format)
	}

	files := c.Arg("files").Strings()

	if ln := len(files); ln > 1 {
		for _, fpath := range files {
			str := fsutil.ReadString(fpath)
			return renderContents(str, fcOpts.format)
		}
	} else if ln == 1 {
		fpath := files[0]
		if fpath[0] != '@' {
			fpath = "@" + fpath
		}
		return renderOneFile(fpath, format)
	}

	// default read stdin
	return renderOneFile(kiteext.DstStdin, format)
}

func renderOneFile(fpath, format string) error {
	sr := kiteext.NewSourceReader(fpath)
	if format == "" && sr.Type() != kiteext.TypeFile {
		format = "markdown" // default as markdown
	}

	str, err := sr.TryReadString()
	if err != nil {
		return err
	}
	return renderContents(str, format)
}

// formatter see like formatters.TTY16m
func renderContents(s, format string) error {
	formatter := "terminal16m"
	if color.IsSupportTrueColor() {
		formatter = "terminal256"
	}
	return quick.Highlight(os.Stdout, s, format, formatter, fcOpts.style)
}

func renderMarkdown(s string) error {
	// 	"github.com/charmbracelet/glamour"
	out, err := glamour.Render(s, fcOpts.style)
	if err == nil {
		stdio.WriteString(out)
	}
	return err
}
