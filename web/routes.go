package web

import "github.com/gookit/rux"

// AddRoutes to rux.Router
func AddRoutes(r *rux.Router) {
	r.StaticDir("/static", "static")

	r.Controller("/", &HomeController{})
}
