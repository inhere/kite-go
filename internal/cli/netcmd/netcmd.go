package netcmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/netcmd/sshcmd"
)

// NetCmd 网络工具命令
var NetCmd = &gcli.Command{
	Name: "net",
	Desc: "Network related commands. ping, telnet ... etc.",
	Subs: []*gcli.Command{
		NewPingCmd(),
		NewNMapCmd(),
		NewNetcatCmd(),
		NewTelnetClientCmd(),
		NewTelnetServerCmd(),
		sshcmd.NewSshExecCmd(),
		sshcmd.NewSshClientCmd(),
	},
}
