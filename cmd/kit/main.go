package main

import (
	"github.com/gookit/goutil/dump"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/app/bootstrap"
)

// dev run:
//	go run ./cmd/kit
//	go run ./cmd/kit
func main() {
	err := bootstrap.Boot(app.App())
	if err != nil {
		panic(err)
	}

	dump.P(app.App().CfgFile(), app.App().Config)
	// do run
	app.Run()
}
