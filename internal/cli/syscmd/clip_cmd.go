package syscmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil/clipboard"
)

// clipCmdOpts struct
type clipCmdOpts struct {
	Verb   bool   `flag:"verbose;display real exec command line;;;v"`
	Read   bool   `flag:"read contents from clipboard, default operate;;;r"`
	Write  string `flag:"write contents to clipboard;false;;w"`
	Output string `flag:"read contents and write to the output;;stdout;o"`
}

// NewClipboardCmd command
func NewClipboardCmd() *gcli.Command {
	var clipOpts = &clipCmdOpts{}

	return &gcli.Command{
		Name:    "clipboard",
		Aliases: []string{"clip"},
		Desc:    "write or read the clipboard contents",
		Config: func(c *gcli.Command) {
			goutil.MustOK(c.UseSimpleRule().FromStruct(clipOpts))
		},
		Func: func(c *gcli.Command, _ []string) error {
			cb := clipboard.New().WithVerbose(clipOpts.Verb)

			if clipOpts.Write != "" {
				_, err := cb.WriteString(clipOpts.Write)
				if err != nil {
					return err
				}
				return cb.Flush()
			}

			// read contents
			str, err := cb.ReadString()
			if err == nil {
				fmt.Println(str)
			}
			return err
		},
	}
}
