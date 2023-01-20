package bootstrap

import (
	"github.com/inhere/kite/app"
	"github.com/inhere/kite/pkg/gitx/gitlab"
)

// addServiceBoot handle
func addServiceBoot(ka *app.KiteApp) {
	ka.AddLoader(app.NewStdLoader(func(ka *app.KiteApp) error {
		glab := gitlab.New()
		err := app.Cfg().MapStruct("gitlab", glab)
		if err != nil {
			return err
		}

		app.Add(app.ObjGlab, glab)
		return nil
	}))
}
