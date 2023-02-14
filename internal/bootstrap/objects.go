package bootstrap

import (
	"github.com/inhere/kite/internal/app"
	"github.com/inhere/kite/pkg/gitx"
	"github.com/inhere/kite/pkg/gitx/github"
	"github.com/inhere/kite/pkg/gitx/gitlab"
	"github.com/inhere/kite/pkg/kiteext"
)

// addServiceBoot handle
func addServiceBoot(ka *app.KiteApp) {
	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		cfg := gitx.NewConfig()
		err := app.Cfg().MapOnExists(app.ObjGit, cfg)
		if err != nil {
			return err
		}

		app.Add(app.ObjGit, cfg)
		return nil
	}, func(ka *app.KiteApp) error {
		// TODO mapstruct dont support set embed struct value.
		cfg := app.Gitx().Clone()
		err := app.Cfg().MapOnExists(app.ObjGlab, cfg)
		if err != nil {
			return err
		}

		glab := gitlab.New(cfg)
		err = app.Cfg().MapOnExists(app.ObjGlab, glab)
		if err != nil {
			return err
		}

		app.Add(app.ObjGlab, glab)
		return nil
	}, func(ka *app.KiteApp) error {
		// TODO mapstruct dont support set embed struct value.
		cfg := app.Gitx().Clone()
		err := app.Cfg().MapOnExists(app.ObjGhub, cfg)
		if err != nil {
			return err
		}

		ghub := github.New(cfg)
		err = app.Cfg().MapOnExists(app.ObjGhub, ghub)
		if err != nil {
			return err
		}

		app.Add(app.ObjGhub, ghub)
		return nil
	})

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

	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		app.OpenMap = app.Cfg().StringMap("quick_open")
		app.PathMap = &kiteext.PathMap{
			Map: app.Cfg().StringMap("pathmap"),
		}
		app.KARun = &kiteext.KiteAliasRun{
			Aliases: app.Cfg().StringMap("aliases"),
		}
		return nil
	})
}
