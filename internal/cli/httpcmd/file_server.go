package httpcmd

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/rux"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/httpserve"
)

// NewFileServerCmd instance
func NewFileServerCmd() *gcli.Command {
	var fsOpts = struct {
		cmdbiz.CommonOpts
		host   string
		port   uint
		prefix string
	}{}

	return &gcli.Command{
		Name:    "fs-server",
		Desc:    "start an simple static file http server",
		Aliases: []string{"fs-serve", "fs-srv", "fss"},
		Config: func(c *gcli.Command) {
			fsOpts.BindWorkdirDryRun(c)
			c.StrOpt2(&fsOpts.host, "host", "custom the file server host", gflag.WithDefault("127.0.0.1"))
			c.UintOpt(&fsOpts.port, "port", "P", 0, "custom the file server port, default will use random `port`")
			c.StrOpt2(&fsOpts.prefix, "prefix,pfx", "custom the static file prefix path for workdir", gflag.WithDefault("/fs"))
		},
		Func: func(c *gcli.Command, args []string) error {
			if fsOpts.port < 1 {
				fsOpts.port = mathutil.SafeUint("1" + timex.Now().DateFormat("md")) // eg: 10425
			}

			srv := httpserve.New(true)
			srv.SetHostPort(fsOpts.host, fsOpts.port)
			wkDir := fsutil.SlashPath(fsOpts.Workdir)
			// wkDir := fsutil.SlashPath(fsutil.ToAbsPath(fsOpts.Workdir))

			r := srv.Rux()
			r.StaticDir(fsOpts.prefix, wkDir)
			colorp.Infof(" - bind URI prefix %s to local path: %s\n", fsOpts.prefix, wkDir)

			cfg := apputil.CmdConfigData2(c)
			if mp := cfg.StrMap("static_dir"); len(mp) > 0 {
				for prefix, dirPath := range mp {
					colorp.Infof("- bind URI prefix %s to dir: %s\n", prefix, dirPath)
					r.StaticDir(prefix, dirPath)
				}
			}

			r.Any("/", func(c *rux.Context) {
				c.JSON(200, rux.M{"msg": "OK", "workdir": wkDir})
			})

			srv.Start()
			return nil
		},
	}
}
