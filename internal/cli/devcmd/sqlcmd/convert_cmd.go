package sqlcmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
)

// Conv2StructCmd convert create table SQL to Go struct
var Conv2StructCmd = &gcli.Command{
	Name:    "struct",
	Aliases: []string{"to-struct", "tostruct", "go-struct"},
	Desc:    "convert create table SQL to Go struct",
	Config: func(c *gcli.Command) {

	},
	Func: func(c *gcli.Command, _ []string) error {
		return errors.New("TODO")
	},
}

// Conv2JSONCmd convert INSERT SQL to JSON object
var Conv2JSONCmd = &gcli.Command{
	Name:    "to-json",
	Aliases: []string{"tomap", "json"},
	Desc:    "convert INSERT ROW SQL to JSON object",
	Config: func(c *gcli.Command) {
		c.AddArg("sql", "the insert SQL. allow: @c")
	},
	Func: func(c *gcli.Command, _ []string) error {
		insertSql, err := apputil.ReadSource(c.Arg("sql").String())
		if err != nil {
			return err
		}

		fieldStr, valueStr := strutil.TrimCut(insertSql, " VALUES ")
		if len(valueStr) == 0 {
			return errors.New("not found keywords 'VALUES' in SQL")
		}

		// spit fields
		fields := strutil.Split(strutil.Trim(fieldStr, " ()"), ",")
		if len(fields) > 0 {
			_, first := strutil.QuietCut(fields[0], "(")
			fields[0] = first
		}

		fields = arrutil.Map(fields, func(obj string) (val string, ok bool) {
			return strings.Trim(obj, " ` "), true
		})
		// dump.P(fields)

		// split values
		values := strutil.Split(strutil.Trim(valueStr, " ();"), ", ")
		anyVals := arrutil.Map(values, func(obj string) (val any, ok bool) {
			obj = strings.TrimSpace(obj)
			if strutil.IsNumeric(obj) {
				return strutil.SafeInt(obj), true
			}
			return strings.Trim(obj, "'\""), true
		})
		// dump.P(values)

		mp := arrutil.CombineToMap(fields, anyVals)
		bs, err := jsonutil.EncodePretty(mp)
		if err != nil {
			return err
		}

		fmt.Println(string(bs))
		return nil
	},
}
