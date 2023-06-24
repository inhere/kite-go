package httpcmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/netutil/httpreq"
	"github.com/gookit/greq"
	"github.com/inhere/kite-go/internal/apputil"
)

var reqOpts = struct {
	// url  string
	data string
	json bool

	method  string
	query   cflag.KVString
	headers cflag.KVString
}{
	query:   cflag.KVString{Sep: "="},
	headers: cflag.KVString{Sep: ":"},
}

// SendRequestCmd instance
var SendRequestCmd = &gcli.Command{
	Name:    "send",
	Aliases: []string{"req", "curl"},
	Desc:    "send http request like curl, ide-http-client",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&reqOpts.json, "j,json", "set use json content type")
		c.StrOpt2(&reqOpts.method, "method, m", "set the reqeust method, default is GET", gflag.WithDefault("GET"))
		c.VarOpt2(&reqOpts.headers, "header, H", "set custom headers, eg: \"Content-Type: application/json\"")
		c.VarOpt2(&reqOpts.query, "query, Q", "append set custom queries, eg: name=inhere")
		c.StrOpt2(&reqOpts.data, "data, d", `set the request body data, eg: '{"name":"inhere"}'`)

		c.AddArg("url", "the url to send request", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		hc := greq.New()
		hc.DefaultMethod(reqOpts.method)

		hc.BeforeSend = func(r *http.Request) {
			cliutil.Yellowln("REQUEST:")
			cliutil.Greenf("%s %s\n\n", r.Method, r.URL.String())
			if len(r.Header) > 0 {
				fmt.Println(httpreq.HeaderToString(r.Header))
			}
		}

		b := hc.Builder()
		b.SetHeaderMap(reqOpts.headers.Data())

		if reqOpts.json {
			b.JSONType()
		}
		if reqOpts.data != "" {
			b.AnyBody(reqOpts.data)
		}
		if !reqOpts.query.IsEmpty() {
			b.WithQuerySMap(reqOpts.query.Data())
		}

		reqOpts := greq.NewOpt()
		apiUrl := c.Arg("url").String()
		resp, err := hc.SendWithOpt(apiUrl, reqOpts)
		if err != nil {
			return err
		}

		cliutil.Yellowln("RESPONSE:")
		if resp.IsEmptyBody() {
			fmt.Print(resp.String())
		} else {
			fmt.Println(resp.String())
		}
		return nil
	},
}

// DecodeQueryCmd instance
var DecodeQueryCmd = &gcli.Command{
	Name:    "dec-query",
	Aliases: []string{"decq", "dq"},
	Desc:    "decode the http query string to structured data",
	Config: func(c *gcli.Command) {
		c.AddArg("query", "the http query string", true)
	},
	Func: func(c *gcli.Command, _ []string) error {
		str, err := apputil.ReadSource(c.Arg("query").String())
		if err != nil {
			return err
		}

		values, err := url.ParseQuery(str)
		if err != nil {
			return err
		}

		mp := make(map[string]any, len(values))
		for key, val := range values {
			if len(val) == 1 {
				mp[key] = val[0]
			} else {
				mp[key] = val
			}
		}

		show.AList("Decoded Query:", mp)
		return nil
	},
}
