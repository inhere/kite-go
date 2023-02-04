package httpserve

import (
	"github.com/gookit/color"
	"github.com/gookit/rux"
)

type HttpServer struct {
	r *rux.Router

	Host string
	Port int
}

func (s HttpServer) Rux() *rux.Router {
	return s.r
}

func New() *HttpServer {
	r := rux.New(rux.EnableCaching)

	// handle error
	r.OnError = func(c *rux.Context) {
		if err := c.FirstError(); err != nil {
			color.Error.Println(err)
			c.HTTPError(err.Error(), 400)
			return
		}
	}

	return &HttpServer{r: r}
}
