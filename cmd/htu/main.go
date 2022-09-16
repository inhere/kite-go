package main

import (
	"github.com/inherelab/kite/internal/cli/doctool"
)

func main() {
	c := doctool.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
