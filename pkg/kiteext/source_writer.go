package kiteext

import (
	"fmt"
	"io"

	"github.com/gookit/goutil/fsutil"
)

// WriterFn type
type WriterFn func(w *SourceWriter)

// SourceWriter struct
type SourceWriter struct {
	err error
	src string
	dst string
	// dstType name
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

func (w *SourceWriter) WriteFrom(r io.Reader) {

}

func (w *SourceWriter) WriteString(s string) (err error) {
	switch w.dstType {
	case TypeFile:
		_, err = fsutil.PutContents(w.dst, s)
	default:
		fmt.Println(s)
	}
	return
}
