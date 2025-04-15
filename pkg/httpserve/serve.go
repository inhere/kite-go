package httpserve

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/rux"
	"github.com/gookit/rux/pkg/handlers"
)

type HttpServer struct {
	r *rux.Router

	Host string
	Port uint
}

func New(printReq bool) *HttpServer {
	rux.Debug(printReq)
	r := rux.New(rux.EnableCaching)
	r.Use(handlers.PanicsHandler())
	if printReq {
		r.Use(handlers.RequestLogger())
	}

	// handle error
	r.OnError = func(c *rux.Context) {
		if err := c.FirstError(); err != nil {
			colorp.Errorln(err)
			c.HTTPError(err.Error(), 400)
			return
		}
	}

	return &HttpServer{r: r}
}

// SetHostPort set host and port
func (s *HttpServer) SetHostPort(host string, port uint) {
	s.Host = host
	s.Port = port
}

// Rux get rux router
func (s *HttpServer) Rux() *rux.Router {
	return s.r
}

// Start http server
func (s *HttpServer) Start() {
	portStr := strutil.SafeString(s.Port)
	// addr := s.Host + ":" + portStr
	s.r.Listen(s.Host, portStr)
}
