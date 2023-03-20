package httpcmd

import (
	"fmt"
	"net/http"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/netutil/httpreq"
	"github.com/gookit/goutil/stdio"
	"github.com/inhere/kite/internal/app"
)

var stOpts = struct {
	envName string
	envFile string
	// topic  string
	domain string
	// ide http client file
	hcFile   string
	tplName  string
	userVars gflag.KVString
	verbose  gcli.VerbLevel
}{
	userVars: cflag.NewKVString(),
}

// SendTemplateCmd instance
var SendTemplateCmd = &gcli.Command{
	Name:    "tpl-send",
	Aliases: []string{"sendtpl", "send-tpl"},
	Desc:    "send http request by a template file or idea http-client file",
	Help: `
## Examples

{$fullCmd} -d gitlab --api api-build.json5 -e prod -v name=order
`,
	Config: func(c *gcli.Command) {
		c.StrOpt2(&stOpts.envName, "env, e", "sets env name for run template")
		c.StrOpt2(&stOpts.envFile, "env-file", "custom sets env file for run template")
		c.StrOpt2(&stOpts.domain, "domain, d", "the domain or topic name")
		c.StrOpt2(&stOpts.hcFile, "http-file, hc-file, hcf", "the ide http client file name or path")
		c.StrOpt2(&stOpts.tplName, "tpl-name, api", "the API template name or file name")
		c.VarOpt2(&stOpts.userVars, "vars, var, v", "custom sets some variables on request. format: `KEY=VALUE`")

		// todo: loop query, send topic, send by template
		// eg:
		// 	kite http send-tpl -d jenkins build --env pre -v name=order
		// 	kite http send-tpl --domain feishu -e dev bot-notify
	},
	Func: func(c *gcli.Command, _ []string) error {
		dc, err := app.HTpl.Domain(stOpts.domain)
		if err != nil {
			return err
		}

		t, err := dc.Lookup(stOpts.hcFile, stOpts.tplName)
		if err != nil {
			return err
		}

		var vs maputil.Data
		vs, err = dc.BuildVars(stOpts.envName, stOpts.envFile)
		if err != nil {
			return err
		}

		uv := stOpts.userVars.Data()
		if len(uv) > 0 {
			vs.LoadSMap(uv)
		}

		if len(vs) > 0 {
			// c.Infof("send request without some variables")
			show.AList("Variables:", vs)
		} else {
			c.Infof("Send request without any variables \n")
		}

		t.BeforeSend = func(r *http.Request) {
			cliutil.Yellowln("REQUEST:")
			cliutil.Greenf("%s %s\n\n", r.Method, r.URL.String())
			fmt.Println(httpreq.HeaderToString(r.Header))
		}

		if err := t.Send(vs, dc.Header); err != nil {
			return err
		}

		cliutil.Yellowln("RESPONSE:")
		if t.Resp.IsEmptyBody() {
			fmt.Print(t.Resp.String())
		} else {
			fmt.Println(t.Resp.String())
		}
		return nil
	},
}

var tiOpts = struct {
	all bool
}{}

// TemplateInfoCmd instance
var TemplateInfoCmd = &gcli.Command{
	Name: "tpl-info",
	// Aliases: []string{"tpl-list"},
	Desc: "list or show loaded config and templates information",
	Config: func(c *gcli.Command) {
		c.BoolOpt2(&tiOpts.all, "list, all, a", "list all configured domains information")

		c.AddArg("domain", "show info for the domain config")
		c.AddArg("group", "show info for the group templates on domain")
		c.AddArg("name", "show template info on the domain.group")
		// todo: loop query, send topic, send by template
		// eg:
		// 	kite http send @jenkins trigger -v env=qa -v name=order
		// 	kite http send @feishu bot-notify
	},
	Func: func(c *gcli.Command, _ []string) error {
		if tiOpts.all {
			c.Infoln("All Domains:")
			dump.NoLoc(app.HTpl.Domains)
			return nil
		}

		domain := c.Arg("domain").String()
		if len(domain) == 0 {
			return errorx.Raw("please input an domain name for show")
		}

		dc, err := app.HTpl.Domain(domain)
		if err != nil {
			return err
		}

		group := c.Arg("group").String()
		if len(group) == 0 {
			show.AList("Domain info", dc)
			return nil
		}

		ts, ok := dc.Templates(group)
		if !ok {
			return errorx.Rawf("the group %q is not fund on domain %q", group, domain)
		}

		name := c.Arg("name").String()
		if len(name) == 0 {
			c.Infoln("Group info:")
			stdio.WriteString(ts.String())
			return nil
		}

		t, err := ts.Lookup(name)
		if err != nil {
			return err
		}

		show.AList("Template info", t)
		return nil
	},
}
