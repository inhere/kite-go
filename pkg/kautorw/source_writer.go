package kautorw

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/stdio"
	"github.com/gookit/goutil/strutil"
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

	openFlush bool
	flushBuf  *bytes.Buffer
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
		dst: strutil.OrElse(dst, "@stdout"),
		// dstType: TypeStdout,
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

// StartFlush start flush buffer contents to dst
func (w *SourceWriter) StartFlush() {
	w.openFlush = true
	w.flushBuf = new(bytes.Buffer)
}

// StopFlush flush all buffer contents to dst
func (w *SourceWriter) StopFlush() error {
	if !w.openFlush {
		return errors.New("not open flush")
	}

	w.openFlush = false
	return w.writeString(w.flushBuf.String())
}

// CleanAndGet clean flush buffer, and get buffer object.
func (w *SourceWriter) CleanAndGet() *bytes.Buffer {
	if !w.openFlush {
		panic("not open flush")
	}

	buf := *w.flushBuf
	w.openFlush = false
	w.flushBuf = nil
	return &buf
}

// WriteFrom write data from reader
func (w *SourceWriter) WriteFrom(r io.Reader) error {
	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return w.WriteString(string(bs))
}

// WriteString write string to dst
func (w *SourceWriter) Write(bs []byte) (err error) {
	if w.openFlush {
		_, err = w.flushBuf.Write(bs)
		return
	}
	return w.writeString(string(bs))
}

// WriteString write string to dst
func (w *SourceWriter) WriteString(s string) (err error) {
	if w.openFlush {
		_, err = w.flushBuf.WriteString(s)
		return
	}
	return w.writeString(s)
}

// write string to dst
func (w *SourceWriter) writeString(s string) (err error) {
	switch w.DstType() {
	case TypeClip:
		err = clipboard.WriteString(s)
	case TypeStdout:
		stdio.WriteString(s)
	default: // write to file
		dst := w.dst
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

// DstType get dst type name
func (w *SourceWriter) DstType() string {
	if w.dstType != "" {
		return w.dstType
	}

	dst := w.dst
	if len(dst) == 0 {
		if w.FallbackType != "" {
			dst = "@" + w.FallbackType
		} else {
			dst = "@stdout"
		}
		w.dst = dst // set to dst
	}

	switch dst {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		w.dstType = TypeClip
	case "@o", "@out", "@stdout", "stdout":
		w.dstType = TypeStdout
	default: // write to file
		w.dstType = TypeFile
	}
	return w.dstType
}
