package fscmd

import (
	"fmt"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/finder"
)

type fileReplaceOpt struct {
	Dir   string      `flag:"desc=the directory for find and replace"`
	Exts  string      `flag:"desc=want render template files exts. multi by comma, like: .go,.md;shorts=ext"`
	Files gcli.String `flag:"desc=the files want replace content, multi by comma"`
	// Include and Exclude.
	// eg: --include name:*.go --exclude name:*_test.go
	Include gcli.Strings `flag:"desc=the files want replace content;shorts=i,add"`
	Exclude gcli.Strings `flag:"desc=the files want replace content;shorts=e,not"`
	// Old from old to new
	Old  string `flag:"desc=the old content want replace"`
	New  string `flag:"desc=the new content replace to"`
	Expr string `flag:"desc=the replace expression. eg: old/new"`
}

// NewReplaceCmd create a command
func NewReplaceCmd() *gcli.Command {
	var opts = fileReplaceOpt{}

	return &gcli.Command{
		Name:    "replace",
		Desc:    "replace content in file(s)",
		Aliases: []string{"re", "rpl", "update"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&opts)
		},
		Func: func(c *gcli.Command, _ []string) error {

			if opts.Files != "" {
				for _, fPath := range opts.Files.Strings() {
					fmt.Println(fPath)
				}
			}

			// find in dir
			if opts.Dir != "" {
				ff := finder.NewFinder(opts.Dir).
					WithExts(strutil.Split(opts.Exts, ","))
				// set finder options
				ff.IncludeRules(opts.Include.Strings())
				ff.ExcludeRules(opts.Exclude.Strings())

				for el := range ff.Find() {
					fmt.Println(el)
				}
			}

			return c.NewErr("TODO")
		},
	}
}
