package kiteext

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
	// emptyAction type on dst is empty.
	emptyAction string
}

// FallbackStdout write to stdout.
func FallbackStdout() ReaderFn {
	return func(sr *SourceReader) {
		sr.emptyAction = TypeStdout
	}
}

func WriteContents(contents, dst string) (string, error) {
	return "", nil
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
	switch dst {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		err = clipboard.WriteString(s)
		w.dstType = TypeClip
	case "", "@o", "@out", "@stdout", "stdout":
		stdio.WriteString(s)
		w.dstType = TypeFile
	default: // to file
		if idx := strings.IndexByte(dst, '@'); idx == 0 {
			dst = dst[1:]
		}
		_, err = fsutil.PutContents(dst, s)
	}
	return
}
