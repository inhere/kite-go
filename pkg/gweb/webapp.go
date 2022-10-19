package gweb

import "github.com/gookit/rux"

// WebApp struct
type WebApp struct {
	router *rux.Router
	srv    *HTTPServer
}

func NewWebApp() *WebApp {
	return &WebApp{}
}
