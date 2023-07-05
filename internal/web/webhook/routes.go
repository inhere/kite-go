package webhook

import (
	"github.com/gookit/rux"
)

// Register routes to rux.Router
func Register(r *rux.Router) {
	r.GET("/", func(c *rux.Context) {
		c.Text(200, "hello, welcome to kite-go webhook server")
	})
	r.Add("/webhook[/{name}]", Webhook, "POST", "GET")
}
