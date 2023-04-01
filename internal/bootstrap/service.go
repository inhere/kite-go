package bootstrap

import (
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/inhere/kite-go/internal/biz/cmdbiz"
	"github.com/inhere/kite-go/pkg/gitx"
	"github.com/inhere/kite-go/pkg/gitx/github"
	"github.com/inhere/kite-go/pkg/gitx/gitlab"
	"github.com/inhere/kite-go/pkg/httptpl"
	"github.com/inhere/kite-go/pkg/kiteext"
	"github.com/inhere/kite-go/pkg/kscript"
	"github.com/inhere/kite-go/pkg/quickjump"
)

// addServiceBoot handle
func addServiceBoot(ka *app.KiteApp) {

	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		app.OpenMap = app.Cfg().StringMap("quick_open")
		app.PathMap = &kiteext.PathMap{
			Aliases: app.Cfg().StringMap("pathmap"),
		}

		app.Vars = kiteext.NewVarMap(app.Cfg().StringMap("global_vars"))
		app.Vars.Aliases["$kite"] = app.Cli().BinName()

		cmdbiz.Kas = app.Cfg().StringMap("aliases")
		return app.Cfg().MapOnExists("proxy_cmd", cmdbiz.ProxyCC)
	})

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
		sr := kscript.NewRunner(func(sr *kscript.Runner) {
			sr.PathResolver = apputil.ResolvePath
		})

		err := app.Cfg().MapOnExists(app.ObjScript, sr)
		if err != nil {
			return err
		}

		// app.Add(app.ObjScript, sr)
		app.Scripts = sr
		return nil
	}, func(ka *app.KiteApp) error {
		plug := &kiteext.PluginRunner{}
		err := app.Cfg().MapOnExists(app.ObjPlugin, plug)
		if err != nil {
			return err
		}

		app.Plugins = plug
		// app.Add(app.ObjPlugin, plug)
		return nil
	}, func(ka *app.KiteApp) error {
		htpl := httptpl.NewManager()
		err := app.Cfg().MapOnExists("http_tpl", htpl)
		if err != nil {
			return err
		}

		htpl.PathResolver = apputil.ResolvePath
		if err := htpl.Init(); err != nil {
			return err
		}

		app.HTpl = htpl
		return nil
	})

	ka.AddBootFuncs(func(ka *app.KiteApp) error {
		app.QJump = quickjump.NewQuickJump()
		app.QJump.PathResolve = apputil.ResolvePath

		err := app.Cfg().MapOnExists("quick_jump", app.QJump)
		if err != nil {
			return err
		}

		return app.QJump.Init()
	})
}
