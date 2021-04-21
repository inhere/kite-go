package main

import (
	"github.com/inherelab/kite/pkg/pacutil"
)

// dev run:
//	go run ./bin/pac
// build run:
//	go build ./bin/pac && ./pac
func main() {
	c := pacutil.PacTools
	c.MustRun(nil)
}
