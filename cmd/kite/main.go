package main

import (
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/bootstrap"
)

// dev run:
//
//	go run ./cmd/kite
//	go run ./cmd/kite <CMD>
//
// install:
//
//	go install ./cmd/kite
func main() {
	// boot and run app
	bootstrap.MustRun(app.App())
}
