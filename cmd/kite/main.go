package main

import (
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/bootstrap"
)

// Dev run:
//
//	go run ./cmd/kite
//	go run ./cmd/kite <CMD>
//
// Debug run:
//	KITE_VERBOSE=debug go run ./cmd/kite <CMD>
//  // Windows PowerShell
//	$env:KITE_VERBOSE="debug"; go run ./cmd/kite <CMD>
//
// Install:
//
//	make install-dev
func main() {
	// boot and run app
	bootstrap.MustRun(app.App())
}
