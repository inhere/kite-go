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
	dst string
	// dstType real dst type name
	dstType string
	srcFile string
	// FallbackType operate type on dst is empty.
	FallbackType string
}

// FallbackStdout write to stdout.
func FallbackStdout() WriterFn {
	return func(sr *SourceWriter) {
		sr.FallbackType = TypeStdout
	}
}

// NewSourceWriter create a new instance
func NewSourceWriter(dst string) *SourceWriter {
	return &SourceWriter{
		dst:     dst,
		dstType: TypeStdout,
	}
}

// WithDst set dst target
func (w *SourceWriter) WithDst(dst string) *SourceWriter {
	w.dst = dst
	return w
}

// SetSrcFile set src file path
func (w *SourceWriter) SetSrcFile(srcPath string) {
	if fsutil.IsFile(srcPath) {
		w.srcFile = srcPath
	}
}

// HasSrcFile check has src file
func (w *SourceWriter) HasSrcFile() bool {
	return w.srcFile != ""
}

func (w *SourceWriter) WriteFrom(r io.Reader) error {
	// TODO
	return nil
}

// WriteString write string to dst
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
	default: // write to file
		if idx := strings.IndexByte(dst, '@'); idx == 0 {
			dst = dst[1:]
		}

		if w.srcFile != "" && dst == "src" {
			dst = w.srcFile
		}

		w.dstType = TypeFile
		_, err = fsutil.PutContents(dst, s)
	}
	return
}
