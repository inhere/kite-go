package main

import (
	"github.com/inhere/kite/internal/cli/doctool"
)

func main() {
	c := doctool.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
