package bootstrap

import (
	"github.com/gookit/gcli/v3"
	"github.com/inhere/kite-go"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli"
)

// BootCli app
func BootCli(_ *app.KiteApp) error {
	cliApp := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Personal developer tool command application"
		a.Version = kite.Version
	})
	// some info
	cliApp.Logo.Text = kite.Banner
	app.Add(app.ObjCli, cliApp)

	// load commands
	cli.Boot(cliApp)
	return nil
}
