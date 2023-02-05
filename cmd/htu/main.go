package main

import (
	"github.com/inhere/kite/internal/cli/toolcmd/doctool"
)

func main() {
	c := doctool.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
