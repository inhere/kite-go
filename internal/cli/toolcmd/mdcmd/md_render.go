package mdcmd

import (
	"github.com/charmbracelet/glamour"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/stdio"
	"github.com/inhere/kite/pkg/kautorw"
)

var mrOpts = struct {
	style string
}{
	style: "auto",
}

// MdRenderCmd instance
var MdRenderCmd = &gcli.Command{
	Name:    "render",
	Aliases: []string{"view", "cat"},
	Desc:    "render markdown file contents on the terminal console",
	Config: func(c *gcli.Command) {
		c.StrOpt2(&mrOpts.style, "style, s", "sets the render style, default is auto.\n allow: auto, dark, dracula, light, notty, pink")
		c.AddArg("file", "want rendered file path", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		fpath := c.Arg("file").String()
		str, err := kautorw.ReadContents("@" + fpath)
		if err != nil {
			return err
		}

		out, err := glamour.Render(str, mrOpts.style)
		if err == nil {
			// fmt.Println(out)
			stdio.WriteString(out)
		}
		return err
	},
}
