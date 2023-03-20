package bootstrap

import (
	"github.com/gookit/gcli/v3"
	kite_go "github.com/inhere/kite-go"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli"
)

// BootCli app
func BootCli(_ *app.KiteApp) error {
	cliApp := gcli.NewApp(func(a *gcli.App) {
		a.Name = "Kite"
		a.Desc = "Personal developer tool command application"
		a.Version = kite_go.Version
	})
	// some info
	cliApp.Logo.Text = kite_go.Banner

	// load commands
	cli.Boot(cliApp)
	app.Add(app.ObjCli, cliApp)
	return nil
}
