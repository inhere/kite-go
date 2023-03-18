package jsoncmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/stdio"
	"github.com/inhere/kite/internal/apputil"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

// JSONToolCmd instance
var JSONToolCmd = &gcli.Command{
	Name: "json",
	Desc: "json tool commands",
	Subs: []*gcli.Command{
		JSONQueryCmd,
		JSONToYAMLCmd,
		JSONFormatCmd,
	},
}

var jvOpts = struct {
	json5 bool
	query string
}{}

// JSONQueryCmd instance
var JSONQueryCmd = &gcli.Command{
	Name:    "view",
	Aliases: []string{"get", "cat", "query"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&jvOpts.json5, "json5, 5", "mark input contents is json5 format")
		c.StrOpt2(&jvOpts.query, "query, path, q, p", "The path for query sub value")

		c.AddArg("json", "input JSON contents for format")
		c.AddArg("path", "The path for query sub value, same of --path")
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(c.Arg("json").String())
		if err != nil {
			return err
		}

		// allow use arg for input path
		if !c.Arg("path").IsEmpty() {
			jvOpts.query = c.Arg("path").String()
		}

		// no query, format and output
		if jvOpts.query == "" {
			return outputFmtJSON(src)
		}

		var mp maputil.Data
		if !jvOpts.json5 {
			// TIP: gjson.Get() cannot find for "dev.host" : {"dev": {"host": "ip:port"}}
			// return outputFmtJSON(gjson.Get(src, jvOpts.query).String())
			if err = json.Unmarshal([]byte(src), &mp); err != nil {
				return err
			}
		} else {
			if err = json5.Unmarshal([]byte(src), &mp); err != nil {
				return err
			}
		}

		// query value
		bs, err := fmtutil.StringOrJSON(mp.Get(jvOpts.query))
		if err != nil {
			return err
		}

		stdio.WriteBytes(bs)
		return nil
	},
}

var jfOpts = struct {
	json5 bool
}{}

// JSONFormatCmd instance
var JSONFormatCmd = &gcli.Command{
	Name:    "format",
	Aliases: []string{"fmt", "pretty"},
	Desc:    "pretty format input JSON(5) contents",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&jfOpts.json5, "json5, 5", "mark input contents is json5 format")

		c.AddArg("json", "input JSON(5) contents for format, allow: @c, @in")
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(c.Arg("json").String())
		if err != nil {
			return err
		}

		return outputFmtJSON(src)
	},
}

func outputFmtJSON(src string) error {
	if len(src) < 24 {
		// only value
		if !strings.ContainsRune(src, ':') {
			stdio.Writeln(src)
			return nil
		}
	}

	var buf bytes.Buffer
	err := json5.Indent(&buf, []byte(src), "", "  ")
	if err != nil {
		return err
	}

	stdio.WriteBytes(buf.Bytes())
	return nil
}

// JSONToYAMLCmd instance
var JSONToYAMLCmd = &gcli.Command{
	Name:    "yaml",
	Aliases: []string{"to-yaml", "to-yml"},
	Desc:    "convert create table SQL to markdown table",
	Config: func(c *gcli.Command) {
		c.AddArg("json", "input JSON contents for convert")
	},
	Func: func(c *gcli.Command, _ []string) error {
		src, err := apputil.ReadSource(c.Arg("json").String())
		if err != nil {
			return err
		}

		bs, err := yaml.JSONToYAML([]byte(src))
		if err != nil {
			return err
		}

		fmt.Println(string(bs))
		return nil
	},
}
