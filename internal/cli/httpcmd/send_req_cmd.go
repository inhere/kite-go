package httpcmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/netutil/httpreq"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/greq"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var reqCmdOpts = struct {
	cmdbiz.CommonOpts
	url string
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
		reqCmdOpts.BindProxyConfirm(c)

		c.BoolOpt2(&reqCmdOpts.json, "j,json", "set use json content type")
		c.StrOpt2(&reqCmdOpts.method, "method, m", "set the request method, default is GET", gflag.WithDefault("GET"))
		c.VarOpt2(&reqCmdOpts.headers, "header, H", `set custom headers, eg: "Content-Type: application/json"`)
		c.VarOpt2(&reqCmdOpts.query, "query, Q", "append set custom queries, eg: name=inhere")
		c.StrOpt2(&reqCmdOpts.data, "data, d", `set the request body data, eg: '{"name":"inhere"}'`)
		c.StrOpt2(&reqCmdOpts.url, "url", "set the request url")

		c.AddArg("url", "set the request url, same of --url")
	},
	Func: func(c *gcli.Command, _ []string) error {
		apiUrl := strutil.OrElse(reqCmdOpts.url, c.Arg("url").String())
		if apiUrl == "" {
			return c.NewErr("the request url is required")
		}

		// create client
		hc := greq.New()

		hc.BeforeSend = func(r *http.Request) error {
			cliutil.Yellowln("REQUEST:")
			cliutil.Greenf("%s %s\n\n", r.Method, r.URL.String())
			if len(r.Header) > 0 {
				fmt.Println(httpreq.HeaderToString(r.Header))
			}
			if reqCmdOpts.data != "" {
				fmt.Println("\n")
				fmt.Println(reqCmdOpts.data)
			}
			return nil
		}

		b := hc.Builder()
		b.SetHeaderMap(reqCmdOpts.headers.Data())

		if reqCmdOpts.json {
			b.JSONType()
			if strings.EqualFold(reqCmdOpts.method, http.MethodGet) {
				b.Method = http.MethodPost
			}
		}
		if reqCmdOpts.data != "" {
			b.AnyBody(reqCmdOpts.data)
		}
		if !reqCmdOpts.query.IsEmpty() {
			b.WithQuerySMap(reqCmdOpts.query.Data())
		}

		request, err := b.Build(reqCmdOpts.method, apiUrl)
		if err != nil {
			return err
		}

		resp, err := hc.SendRequest(request)
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

		show.AList("Decoded Query:", mp, func(opts *show.ListOption) {
			// opts.KeyStyle = "" // disable color for key
		})
		return nil
	},
}
