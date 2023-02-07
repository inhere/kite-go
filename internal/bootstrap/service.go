package bootstrap

import (
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/gitlab"
	"github.com/inhere/kite/pkg/kiteext"
)

// addServiceBoot handle
func addServiceBoot(ka *app.KiteApp) {
	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		cfg := gitx.NewConfig()
		err := app.Cfg().MapOnExists(app.ObjGitx, cfg)
		if err != nil {
			return err
		}

		app.Add(app.ObjGitx, cfg)
		return nil
	}).
		AddLoader(app.NewStdLoader(func(ka *app.KiteApp) error {
			glab := gitlab.New()
			err := app.Cfg().MapOnExists(app.ObjGlab, glab)
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
