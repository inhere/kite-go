package main

import (
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/app/bootstrap"
)

// dev run:
//
//	go run ./cmd/kit
//	go run ./cmd/kit -h
//
// install:
//
//	go install ./cmd/kit
func main() {
	bootstrap.MustBoot(app.App())

	// do run
	app.Run()
}
