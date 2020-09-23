package httpserve

import "github.com/gookit/rux"

func NewServe() *rux.Router {
	return rux.New(rux.EnableCaching)
}
