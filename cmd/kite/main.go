package main

import (
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/internal/bootstrap"
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
