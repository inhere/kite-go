package netcmd

import "github.com/gookit/gcli/v3"

// NewTelnetServerCmd creates a new TelnetServerCmd
func NewTelnetServerCmd() *gcli.Command {
	tsOpts := struct {
		host string
		port int
	}{}

	return &gcli.Command{
		Name:    "telnet-server",
		Desc:    "start a telnet server",
		Aliases: []string{"ts", "telnet-s"},
		Config: func(c *gcli.Command) {
			c.StrOpt(&tsOpts.host, "host", "H", "127.0.0.1", "telnet server host")
			c.IntOpt(&tsOpts.port, "port", "p", 23, "telnet server port")
		},
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},
	}
}
