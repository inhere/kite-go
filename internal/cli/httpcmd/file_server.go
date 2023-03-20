package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/rux"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
)

// NewFileServerCmd instance
func NewFileServerCmd() *gcli.Command {
	var fsOpts = struct {
		cmdbiz.CommonOpts
		port uint
	}{}

	return &gcli.Command{
		Name:    "fs-server",
		Desc:    "start an simple file http server",
		Aliases: []string{"fs-serve", "fs-srv", "fss"},
		Config: func(c *gcli.Command) {
			fsOpts.BindCommonFlags1(c)
			c.UintOpt(&fsOpts.port, "port", "P", 0, "custom the echo server port, default will use random `port`")
		},
		Func: func(c *gcli.Command, args []string) error {
			if fsOpts.port < 1 {
				fsOpts.port = uint(mathutil.RandInt(6000, 9999))
			}

			srv := rux.New(func(r *rux.Router) {})
			srv.StaticDir("/fs", fsOpts.Workdir)

			cfg := apputil.CmdConfigData2(c)
			if mp := cfg.StrMap("static_dir"); len(mp) > 0 {
				for prefix, dirPath := range mp {
					srv.StaticDir(prefix, dirPath)
				}
			}

			srv.Any("/", func(c *rux.Context) {
				c.JSON(200, rux.M{"msg": "OK"})
			})

			srv.Listen(mathutil.String(fsOpts.port))
			return srv.Err()
		},
	}
}
