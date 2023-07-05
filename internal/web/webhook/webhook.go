package webhook

import "github.com/gookit/rux"

// Webhook api for the application
func Webhook(c *rux.Context) {
	c.JSON(200, rux.M{"hello": "welcome"})
}
