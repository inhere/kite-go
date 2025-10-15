package main

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go/internal/cli/devcmd/envcmd"
)

// main ktenv 程序入口
//
// Dev run:
//
//	go run ./cmd/ktenv
//	go run ./cmd/ktenv <CMD>
//
// Debug run:
//	KITE_VERBOSE=debug go run ./cmd/ktenv <CMD>
//  // Windows PowerShell
//	$env:KITE_VERBOSE="debug" go run ./cmd/ktenv <CMD>
//
func main() {
	em := NewEnvManageCmd()
	em.MustRun(nil)
}

func NewEnvManageCmd() *gcli.Command {
	return &gcli.Command{
		Name: "kenv",
		// Aliases: []string{"envs", "useenv"},
		Desc: "manage local development environment SDK",
		Help: `
Commands Usage:
  use <sdk:version>...     Activate SDK versions
    -s, --save             Save configuration to project file
  unuse <sdk>...           Deactivate SDKs
  add <sdk:version>...     Download and install SDK versions
  list [sdk]               List installed SDKs

Examples:
  ktenv use node:18 go:1.21
  ktenv use -s node:lts
  ktenv unuse node
  ktenv add go:1.22
  ktenv list
  ktenv list go

Supported SDKs:
  go, node, java, flutter

Version formats:
  <sdk>:<version>         Exact version (go:1.21.5)
  <sdk>:<major>           Latest patch version (node:18)
  <sdk>:lts               Long-term support version
  <sdk>:latest            Latest stable version
`,
		Subs: []*gcli.Command{
			envcmd.NewEnvListCmd(),
			envcmd.NewEnvAddCmd(),
			envcmd.NewEnvRemoveCmd(),
			envcmd.NewEnvUseCmd(),
			envcmd.NewEnvShellCmd(),
			envcmd.NewEnvConfigCmd(),
		},
	}
}
