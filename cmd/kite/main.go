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
	bootstrap.MustBoot(app.App())

	// do run
	app.Run()
}
