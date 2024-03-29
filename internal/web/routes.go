package web

import (
	"github.com/gookit/rux"
	"github.com/inhere/kite-go/internal/web/controller"
)

// AddRoutes to rux.Router
func AddRoutes(r *rux.Router) {
	r.StaticDir("/static", "static")

	r.Controller("/", &controller.HomeController{})
	r.Controller("/tasks", &controller.TaskController{})
	r.Controller("/json", &controller.JSONPage{})
}
