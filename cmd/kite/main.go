package main

import (
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/app/bootstrap"
)

// dev run:
//
//	go run ./bin/kite
//	go run ./bin/kite
//
// install:
//
//	go install ./cmd/kite
func main() {
	// boot and run app
	bootstrap.MustRun(app.App())
}
