package main

import "github.com/inhere/kite/internal/cli/toolcmd/doccmd"

func main() {
	c := doccmd.DocumentCmd
	c.Name = "htu"

	c.MustRun(nil)
}
