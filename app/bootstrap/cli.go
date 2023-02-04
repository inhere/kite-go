package bootstrap

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite"
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/internal/cli"
)

// BootCli app
func BootCli(_ *app.KiteApp) error {
	cliApp := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Kite CLI tool application"
		a.Version = kite.Version
	})
	// some info
	cliApp.Logo.Text = kite.Banner

	// load commands
	cli.Boot(cliApp)
	app.Add(app.ObjCli, cliApp)
	return nil
}
