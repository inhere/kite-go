package syscmd

import (
	"path/filepath"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/strutil"
)

var eiOpts = struct {
	search  string
	inValue bool
	expand  bool
}{}

// NewEnvInfoCmd instance
func NewEnvInfoCmd() *gcli.Command {
	return &gcli.Command{
		Name: "env",
		Desc: "display system env information",
		// Aliases: []string{"exec", "script"},
		Config: func(c *gcli.Command) {
			c.StrOpt2(&eiOpts.search, "search,s", "The keywords for search ENV information")
			c.BoolOpt2(&eiOpts.inValue, "value,v", "Match ENV value on search. default only match key.")
			c.BoolOpt2(&eiOpts.expand, "expand,e", "expand ENV value for give name, useful for PATH.")

			c.AddArg("name", "display ENV value by the name")
		},
		Func: envInfoHandle,
	}
}

func envInfoHandle(c *gcli.Command, _ []string) error {
	if eiOpts.search != "" {
		founded := make(map[string]string)
		for name, val := range envutil.Environ() {
			if strutil.IContains(name, eiOpts.search) {
				founded[name] = val
			} else if eiOpts.inValue && strutil.IContains(val, eiOpts.search) {
				founded[name] = val
			}
		}

		show.AList("Search Result", founded)
		return nil
	}

	getName := c.Arg("name").String()
	if getName == "" {
		show.AList("ENV Information", envutil.Environ())
		return nil
	}

	for name, val := range envutil.Environ() {
		if strutil.IEqual(name, getName) {
			if eiOpts.expand {
				show.AList("Expand "+name, filepath.SplitList(val))
			} else {
				c.Println(val)
			}
			break
		}
	}
	return nil
}
