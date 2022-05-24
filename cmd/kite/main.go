package main

import (
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/app/bootstrap"
)

// dev run:
//	go run ./bin/kite
//	go run ./bin/kite
func main() {
	err := bootstrap.Boot(app.App())
	if err != nil {
		panic(err)
	}

	// do run
	app.Run()
}
