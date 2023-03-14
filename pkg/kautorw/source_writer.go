package kautorw

import (
	"io"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/stdio"
	"github.com/gookit/goutil/sysutil/clipboard"
)

// WriterFn type
type WriterFn func(w *SourceWriter)

// SourceWriter struct
type SourceWriter struct {
	err error
	src string
	dst string
	// dstType real dst type name
	dstType string
	// FallbackType operate type on dst is empty.
	FallbackType string
}

// FallbackStdout write to stdout.
func FallbackStdout() WriterFn {
	return func(sr *SourceWriter) {
		sr.FallbackType = TypeStdout
	}
}

func NewSourceWriter(dst string) *SourceWriter {
	return &SourceWriter{
		dst:     dst,
		dstType: TypeStdout,
	}
}

func (w *SourceWriter) WriteFrom(r io.Reader) error {

	return nil
}

func (w *SourceWriter) WriteString(s string) (err error) {
	dst := w.dst
	if len(dst) == 0 && w.FallbackType != "" {
		dst = "@" + w.FallbackType
	}

	switch dst {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		w.dstType = TypeClip
		err = clipboard.WriteString(s)
	case "@o", "@out", "@stdout", "stdout":
		w.dstType = TypeStdout
		stdio.WriteString(s)
	default: // to file
		if idx := strings.IndexByte(dst, '@'); idx == 0 {
			dst = dst[1:]
		}
		w.dstType = TypeFile
		_, err = fsutil.PutContents(dst, s)
	}
	return
}
