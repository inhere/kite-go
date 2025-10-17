package main

import (
	"github.com/inhere/kite-go/internal/cli/xenvcmd"
)

// main xenv 程序入口
//
// Dev run:
//
//	go run ./cmd/xenv
//	go run ./cmd/xenv <CMD>
//
// Debug run:
//	KITE_VERBOSE=debug go run ./cmd/xenv <CMD>
//  // Windows PowerShell
//	$env:KITE_VERBOSE="debug" go run ./cmd/xenv <CMD>
//
func main() {
	xe := xenvcmd.XEnvCmd
	xe.Help = `
Commands Usage:
  use <sdk:version>...     Activate SDK versions
    -s, --save             Save configuration to project file
  unuse <sdk>...           Deactivate SDKs
  add <sdk:version>...     Download and install SDK versions
  list [sdk]               List installed SDKs

Examples:
  xenv use node:18 go:1.21
  xenv use -s node:lts
  xenv unuse node
  xenv add go:1.22
  xenv list
  xenv list go

Supported SDKs:
  go, node, java, flutter

Version formats:
  <sdk>:<version>         Exact version (go:1.21.5)
  <sdk>:<major>           Latest patch version (node:18)
  <sdk>:lts               Long-term support version
  <sdk>:latest            Latest stable version
`
	xe.MustRun(nil)
}
