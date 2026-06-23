package kserve

import "github.com/gookit/rux/v2"

// WebApp struct
type WebApp struct {
	router *rux.Router
	srv    *HTTPServer
}

func NewWebApp() *WebApp {
	return &WebApp{}
}
