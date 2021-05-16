package web

import (
	"strings"

	"github.com/gookit/rux"
)

type TaskController struct {}

// AddRoutes to rux.Router
func (c *TaskController) AddRoutes(r *rux.Router) {
	r.GET("", c.Index)
	r.GET("add", c.Add)
}

// Index page for the application
func (*TaskController) Index(c *rux.Context) {
	c.JSON(200, rux.M{"hello": "welcome"})
}

// Add page for the application
func (*TaskController) Add(c *rux.Context) {
	jobId := c.Query("jobId")
	jobId = strings.TrimSpace(jobId)
	if jobId == "" {
		c.AbortThen().JSON(406, rux.M{"error": "invalid job ID"})
	}

	c.JSON(200, rux.M{"hello": "add task"})
}
