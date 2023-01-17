package main

import (
	"github.com/inhere/kite/pkg/pacutil"
)

// dev run:
//
//	go run ./bin/pacgo
//
// build run:
//
//	go build ./bin/pacgo && ./pac
func main() {
	c := pacutil.PacTools
	c.MustRun(nil)
}
