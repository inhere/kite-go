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
	"github.com/gookit/goutil/netutil/httpctype"
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
		hc.Method(reqOpts.method)

		if reqOpts.json {
			hc.ContentType(httpctype.JSON)
		}
		if reqOpts.data != "" {
			hc.AnyBody(reqOpts.data)
		}

		hc.BeforeSend = func(r *http.Request) {
			cliutil.Yellowln("REQUEST:")
			cliutil.Greenf("%s %s\n\n", r.Method, r.URL.String())
			if len(r.Header) > 0 {
				fmt.Println(httpreq.HeaderToString(r.Header))
			}
		}

		apiUrl := c.Arg("url").String()
		if !reqOpts.query.IsEmpty() {
			apiUrl = httpreq.AppendQueryToURLString(apiUrl, httpreq.ToQueryValues(reqOpts.query.Data()))
		}

		reqOpts := &greq.Option{
			HeaderM: reqOpts.headers.Data(),
		}
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

		show.AList("Decoded Query:", values)
		return nil
	},
}
