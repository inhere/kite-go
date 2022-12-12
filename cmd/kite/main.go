package main

import (
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/app/bootstrap"
)

// dev run:
//
//	go run ./bin/kite
//	go run ./bin/kite
func main() {
	bootstrap.MustBoot(app.App())

	// do run
	app.Run()
}
