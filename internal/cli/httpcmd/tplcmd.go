package httpcmd

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/netutil/httpreq"
	"github.com/gookit/goutil/stdio"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

var stOpts = struct {
	cmdbiz.CommonOpts
	envName string
	envFile string
	// topic  string
	domain string
	// ide http client file
	hcFile   string
	tplName  string
	timeout  int // ms
	userVars gflag.KVString
	verbose  bool
	plugins  gflag.String
}{
	timeout:  500,
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

## use variable in template

{$fullCmd} -d gitlab --plug git,fs --api api-build.json5 -e prod -v group={{git.group}} -v repoName={{git.repo}}
`,
	Config: func(c *gcli.Command) {
		stOpts.BindCommonFlags(c)
		c.IntOpt2(&stOpts.timeout, "timeout, t", "sets the request timeout, unit: ms. default: 500")
		c.StrOpt2(&stOpts.envName, "env, e", "sets env name for run template")
		c.StrOpt2(&stOpts.envFile, "env-file", "custom sets env file for run template")
		c.StrOpt2(&stOpts.domain, "domain, d", "the domain or topic name")
		c.StrOpt2(&stOpts.hcFile, "http-file, hc-file, hcf", "the ide http client file name or path")
		c.StrOpt2(&stOpts.tplName, "tpl-name, api", "the API template name or file name")
		c.VarOpt2(&stOpts.plugins, "plugin,plug", "enable some plugins on exec request. allow:git,fs\ne.g. --plugin=plugin1,plugin2")
		c.VarOpt2(&stOpts.userVars, "vars, var, v", "custom sets some variables on request. format: `KEY=VALUE`")
		c.BoolOpt2(&stOpts.verbose, "verbose, vv", `show more info about request and response`)

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

		vs.LoadSMap(stOpts.userVars.Data())

		if len(stOpts.plugins) > 0 {
			wDir := c.WorkDir()
			names := stOpts.plugins.Strings()
			for _, name := range names {
				switch name {
				case "fs":
					vs.LoadSMap(map[string]string{
						"fs.dir":  fsutil.Name(wDir),
						"fs.path": wDir,
					})
				case "git":
					if gitw.IsGitDir(wDir) {
						lp := gitw.NewRepo(wDir)
						if ri := lp.FirstRemoteInfo(); ri != nil {
							vs.LoadSMap(map[string]string{
								"git.group":    ri.Group,
								"git.repo":     ri.Repo,
								"git.repoPath": ri.RepoPath(),
							})
						}
					}
				}
			}
		}

		t.SetTimeout(stOpts.timeout)

		if stOpts.verbose {
			show.AList("Request Options:", map[string]any{
				"timeout(ms)": t.Timeout,
			})

			if len(vs) > 0 {
				// c.Infof("send request without some variables")
				show.AList("Variables:", vs)
			} else {
				c.Infoln("Send template request without any variables")
			}

			t.BeforeSend = func(r *http.Request, b *bytes.Buffer) {
				cliutil.Yellowln("REQUEST:")
				cliutil.Greenf("%s %s\n\n", r.Method, r.URL.String())
				if len(r.Header) > 0 {
					fmt.Println(httpreq.HeaderToString(r.Header))
				}

				if b != nil && b.Len() > 0 {
					fmt.Println(b.String())
				}
				colorp.Cyanln("\n-------------------------------------------------------------------------\n", "")
			}
			t.AfterSend = func(resp *httpreq.Resp, err error) {
				if err != nil {
					return
				}

				cliutil.Yellowln("RESPONSE:")
				if resp.IsEmptyBody() {
					fmt.Print(resp.String())
				} else {
					fmt.Println(resp.String())
				}
			}
		}

		opt := httpreq.NewOpt()
		if err = t.Send(vs, dc.Header, opt); err != nil {
			return err
		}

		if !stOpts.verbose {
			fmt.Println(t.Resp.BodyString())
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
