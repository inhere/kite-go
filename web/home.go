package web

import (
	"os"

	"github.com/gookit/kite/app"
	"github.com/gookit/rux"
	"github.com/gookit/view"
)

// HomeController struct
type HomeController struct{}

// AddRoutes to rux.Router
func (c *HomeController) AddRoutes(r *rux.Router) {
	r.GET("", c.Index)
	r.GET("apidoc", c.ApiDoc)
	r.GET("about[.html]", c.About)
}

// Index page for the application
func (*HomeController) Index(c *rux.Context) {
	c.JSON(200, rux.M{"hello": "welcome"})
}

// About page for the application
func (*HomeController) About(c *rux.Context) {
	c.JSON(200, rux.M{"hello": "welcome"})
}

// ApiDoc page for display swagger doc
func (*HomeController) ApiDoc(c *rux.Context) {
	swagFile := "static/apidoc/swagger.json"
	fInfo, err := os.Stat(swagFile)
	if err != nil {
		c.AbortWithStatus(404, "swagger doc file not exists")
		return
	}

	data := map[string]string{
		"EnvName":    "prod",
		"AppName":    "Kite Application",
		"JsonFile":   "/" + swagFile,
		"SwgUIPath":  "/static/swaggerui",
		"AssetPath":  "/static",
		"UpdateTime": fInfo.ModTime().Format(app.DateFormat),
	}

	// c.HTML(200, nil)
	c.AddError(view.Partial(c.Resp, "swagger.tpl", data))
}
