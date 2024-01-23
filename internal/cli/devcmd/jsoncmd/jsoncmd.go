package jsoncmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/stdio"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

// JSONToolCmd instance
var JSONToolCmd = &gcli.Command{
	Name: "json",
	Desc: "json format, convert tool commands",
	Subs: []*gcli.Command{
		JSONQueryCmd,
		JSONToYAMLCmd,
		JSONFormatCmd,
	},
}

var jvOpts = struct {
	json5 bool
	query string
	// compressed, not format output
	compressed bool
}{}

// JSONQueryCmd instance
var JSONQueryCmd = &gcli.Command{
	Name:    "view",
	Aliases: []string{"get", "cat", "query"},
	Desc:    "format and query value from JSON(5) contents",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&jvOpts.json5, "json5, 5", "mark input contents is json5 format")
		c.StrOpt2(&jvOpts.query, "query, path, q, p", "The path for query sub value")
		c.BoolOpt2(&jvOpts.compressed, "compressed, c", "compressed output, not format")

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
			if err = json.Unmarshal([]byte(src), &mp); err != nil {
				return err
			}
		} else {
			if err = json5.Unmarshal([]byte(src), &mp); err != nil {
				return err
			}
		}

		// query value
		value := mp.Get(jvOpts.query)
		s, err := strutil.ToStringWith(value)
		if err == nil {
			stdio.Writeln(s)
			return nil
		}

		// err != nil: use json format output
		bs, err1 := json.Marshal(value)
		if err1 != nil {
			return err1
		}

		if jvOpts.compressed {
			stdio.WritelnBytes(bs)
			return nil
		}

		// format output
		var buf bytes.Buffer
		err = json.Indent(&buf, bs, "", "    ")
		if err != nil {
			return err
		}

		stdio.WriteBytes(buf.Bytes())
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
