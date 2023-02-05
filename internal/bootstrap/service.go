package bootstrap

import (
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/pkg/gitx/gitlab"
	"github.com/inhere/kite/pkg/kiteext"
)

// addServiceBoot handle
func addServiceBoot(ka *app.KiteApp) {
	ka.AddLoader(app.NewStdLoader(func(ka *app.KiteApp) error {
		glab := gitlab.New()
		err := app.Cfg().MapOnExists("gitlab", glab)
		if err != nil {
			return err
		}

		app.Add(app.ObjGlab, glab)
		return nil
	}))

	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		sr := &kiteext.ScriptRunner{}
		err := app.Cfg().MapOnExists(app.ObjScript, sr)
		if err != nil {
			return err
		}

		app.Add(app.ObjScript, sr)
		return nil
	}, func(ka *app.KiteApp) error {
		plug := &kiteext.PluginRunner{}
		err := app.Cfg().MapOnExists(app.ObjPlugin, plug)
		if err != nil {
			return err
		}

		app.Add(app.ObjPlugin, plug)
		return nil
	})
}
