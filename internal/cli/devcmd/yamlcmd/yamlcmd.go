package jsoncmd

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/errorx"
	"github.com/inhere/kite-go/internal/apputil"
)

// YamlToolCmd instance
var YamlToolCmd = &gcli.Command{
	Name: "yaml",
	Desc: "yaml format contents tool commands",
	Subs: []*gcli.Command{
		YamlViewCmd,
		YamlCheckCmd,
		YamlToJSONCmd,
	},
}

// YamlViewCmd instance
var YamlViewCmd = &gcli.Command{
	Name:    "view",
	Aliases: []string{"cat", "fmt"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}

// YamlToJSONCmd instance
var YamlToJSONCmd = &gcli.Command{
	Name:    "json",
	Aliases: []string{"to-json"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {
		c.AddArg("yaml", "input yaml contents for convert")
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(c.Arg("yaml").String())
		if err != nil {
			return err
		}

		bs, err := yaml.YAMLToJSON([]byte(src))
		if err != nil {
			return err
		}

		fmt.Println(string(bs))
		return nil
	},
}

// YamlCheckCmd instance
var YamlCheckCmd = &gcli.Command{
	Name:    "check",
	Aliases: []string{"validate"},
	Desc:    "validate YAML contents format",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errorx.New("TODO")
	},
}
