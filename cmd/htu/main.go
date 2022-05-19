package main

import (
	"github.com/inherelab/kite/internal/cmd/doctool"
)

func main() {
	c := doctool.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
