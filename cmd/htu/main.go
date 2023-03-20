package main

import "github.com/inhere/kite-go/internal/cli/toolcmd/doccmd"

func main() {
	c := doccmd.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
