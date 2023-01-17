package bootstrap

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/cli"
)

func BootCli(ka *app.KiteApp) error {
	cliApp := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"

		a.Version = kite.Version
	})

	// load commands
	cli.Boot(cliApp)
	app.Add("cli", cliApp)
	return nil
}
