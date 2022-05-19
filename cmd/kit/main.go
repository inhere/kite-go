package main

import (
	"github.com/inherelab/kite"
	"github.com/inherelab/kite/app"
	"github.com/inherelab/kite/internal/cmd"
)

// dev run:
//	go run ./cmd/kit
//	go run ./cmd/kit
func main() {
	cli := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"

		a.Version = kite.Version
	})

	app.Boot(cli)

	// load commands
	cmd.Boot(cli)

	// do run
	cli.Run(nil)
}
