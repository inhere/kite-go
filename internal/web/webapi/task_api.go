package webapi

import "github.com/gookit/rux/v2"

// TaskApiController struct
type TaskApiController struct{}

// Index api for the application
func (*TaskApiController) Index(c *rux.Context) {
	c.JSON(200, rux.M{"hello": "welcome"})
}
