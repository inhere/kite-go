package httpcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/timex"
	"github.com/inhere/kite-go/internal/web/webhook"
	"github.com/inhere/kite-go/pkg/httpserve"
)

var hookSrvOpts = struct {
	Port   uint   `flag:"desc=custom the webhook server port, default will use random port;shorts=P"`
	Config string `flag:"desc=custom the webhook server config file path;shorts=C"`
	Debug  bool   `flag:"desc=enable debug mode;shorts=D"`
}{}

// NewHookServerCmd new command
func NewHookServerCmd() *gcli.Command {
	return &gcli.Command{
		Name: "hook-serve",
		Desc: "start a http server for receive webhook request",
		Help: `
## Examples

  # start a webhook server
  {$fullCmd} -p 8080
  # start a webhook server with config file
  {$fullCmd} -c ./hook-server.yml

  # access the server
  curl -X POST http://localhost:8080/webhook -d '{"name": "test"}'
`,
		Aliases: []string{"webhook", "hook-server"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&hookSrvOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			s := httpserve.New(hookSrvOpts.Debug)

			if hookSrvOpts.Port < 500 {
				hookSrvOpts.Port = mathutil.SafeUint("1" + timex.Now().DateFormat("md")) // eg: 10425
			}
			s.SetHostPort("127.0.0.1", hookSrvOpts.Port)

			webhook.Register(s.Rux())

			// quick start
			// r.Listen("127.0.0.1", mathutil.String(hookSrvOpts.Port))
			s.Start()
			return nil
		},
	}
}
