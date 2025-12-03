package netcmd

import "github.com/gookit/gcli/v3"

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
	},
}
