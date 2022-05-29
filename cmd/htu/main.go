package main

import (
	"github.com/inherelab/kite/internal/command/doctool"
)

func main() {
	c := doctool.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
