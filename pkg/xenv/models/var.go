package models

import (
	"fmt"

	"github.com/gookit/goutil/x/ccolor"
)

var (
	DebugMode bool
)

// Debugf prints debug messages
func Debugf(format string, args ...any) {
	if DebugMode {
		ccolor.Printf("<cyan>DEBUG</>: "+format, args...)
	}
}

func Debugln(args ...any) {
	if DebugMode {
		ccolor.Println("<cyan>DEBUG</>: ", fmt.Sprint(args...))
	}
}

