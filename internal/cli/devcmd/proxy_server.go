package devcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/pkg/proxysrv"
)

type pSrvOptions struct {
	Port   int    `flag:"desc=proxy server port;default=8080"`
	Rules  string `flag:"desc=proxy rules file path;default=proxy-rules.txt"`
	Config string `flag:"desc=proxy config file path;default=.dev-proxy.yml"`
}

// NewProxyServerCmd create a new command.
func NewProxyServerCmd() *gcli.Command {
	psOpts := pSrvOptions{}

	return &gcli.Command{
		Name:    "proxy-server",
		Desc:    "Start a http proxy server for development",
		Aliases: []string{"proxy-s", "proxy-srv"},
		Config: func(c *gcli.Command) {
			c.MustFromStruct(&psOpts)
		},
		Func: func(c *gcli.Command, args []string) error {
			ps := proxysrv.NewProxySrv()
			// TODO
			return ps.Start()
		},
	}
}
