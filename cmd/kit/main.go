package main

import (
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/app/bootstrap"
)

// dev run:
//
//	go run ./cmd/kit
//	go run ./cmd/kit -h
func main() {
	bootstrap.MustBoot(app.App())

	// do run
	app.Run()
}
