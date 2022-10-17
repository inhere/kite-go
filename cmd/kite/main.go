package main

import (
	"github.com/gookit/goutil"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/app/bootstrap"
)

// dev run:
//
//	go run ./bin/kite
//	go run ./bin/kite
func main() {
	goutil.MustOK(bootstrap.Boot(app.App()))

	// do run
	app.Run()
}
