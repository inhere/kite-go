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
	err := bootstrap.Boot(app.App())
	if err != nil {
		panic(err)
	}

	// do run
	app.Run()
}
