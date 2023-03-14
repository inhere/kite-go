package kautorw

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil/clipboard"
)

// ReaderFn type
type ReaderFn func(r *SourceReader)

// SourceReader struct
type SourceReader struct {
	buf bytes.Buffer
	err error
	src string
	// src real type
	typ string
	// has read contents by src
	hasRead bool
	// TypeOnEmpty operate type on src is empty.
	TypeOnEmpty string
	// DefaultAsFile type on src type not match.
	DefaultAsFile bool
	// TrimSpace for read contents
	TrimSpace bool
	// CheckResult return error on result is empty
	CheckResult bool
}

// WithTypeOnEmpty setting.
func WithTypeOnEmpty(typ string) ReaderFn {
	return func(sr *SourceReader) { sr.TypeOnEmpty = typ }
}

// TryStdinOnEmpty try read from stdin.
func TryStdinOnEmpty() ReaderFn {
	return func(sr *SourceReader) { sr.TypeOnEmpty = TypeStdin }
}

// WithTrimSpace on read contents.
func WithTrimSpace() ReaderFn {
	return func(sr *SourceReader) { sr.TrimSpace = true }
}

// WithDefaultAsFile on read contents.
func WithDefaultAsFile() ReaderFn {
	return func(sr *SourceReader) { sr.DefaultAsFile = true }
}

// WithCheckResult on read contents.
func WithCheckResult() ReaderFn {
	return func(sr *SourceReader) { sr.CheckResult = true }
}

// NewSourceReader instance
func NewSourceReader(src string, fns ...ReaderFn) *SourceReader {
	sr := &SourceReader{
		src: src,
	}

	sr.CheckResult = true
	for _, fn := range fns {
		fn(sr)
	}
	return sr
}

// TryString return string and error
func (r *SourceReader) TryString() (string, error) {
	return r.ReadString(), r.err
}

// Buffer get
func (r *SourceReader) Buffer() *bytes.Buffer {
	return &r.buf
}

// Reset buffer
func (r *SourceReader) Reset() *SourceReader {
	if r.buf.Len() > 0 {
		r.buf.Reset()
	}

	r.hasRead = false
	return r
}

// ReadClip return string and error
func (r *SourceReader) ReadClip() *SourceReader {
	r.tryReadClip()
	return r
}

// tryReadClip return string and error
func (r *SourceReader) tryReadClip() {
	ln := len(r.src)
	if ln > 12 {
		r.directToBuf()
		return
	}

	switch r.src {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		r.typ = TypeClip
		r.err = clipboard.Std().ReadTo(&r.buf)
		r.hasRead = true
	default:
		r.directToBuf()
	}
}

// ReadStdin handle
func (r *SourceReader) ReadStdin() *SourceReader {
	r.tryReadStdin()
	return r
}

// tryReadStdin return string and error
func (r *SourceReader) tryReadStdin() {
	ln := len(r.src)
	if ln > 8 {
		r.directToBuf()
		return
	}

	switch r.src {
	case "@in", "@stdin", "stdin":
		r.typ = TypeStdin
		_, r.err = r.buf.ReadFrom(os.Stdin)
		r.hasRead = true
	default:
		r.directToBuf()
	}
}

// TryReadString return string and error
func (r *SourceReader) TryReadString() (string, error) {
	return r.ReadString(), r.err
}

var emptyResultErr = errors.New("input is empty")

// ReadString return string
func (r *SourceReader) ReadString() string {
	if r.buf.Cap() == 0 {
		r.tryReadAny()
	}
	defer r.buf.Reset()

	s := r.buf.String()
	if r.TrimSpace {
		s = strings.TrimSpace(s)
	}

	if r.CheckResult && len(s) == 0 {
		r.err = emptyResultErr
	}
	return s
}

// tryReadAny return string
func (r *SourceReader) tryReadAny() {
	ln := len(r.src)
	if ln > 86 {
		r.directToBuf()
		return
	}

	src := r.src
	if r.TypeOnEmpty != "" && ln == 0 {
		src = "@" + r.TypeOnEmpty
	}

	r.tryReadSrc(src)
}

// HasRead bool
func (r *SourceReader) tryReadSrc(src string) {
	switch src {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		r.typ = TypeClip
		r.err = clipboard.Std().ReadTo(&r.buf)
		r.hasRead = true
	case "@i", "@in", "@stdin", "stdin":
		r.typ = TypeStdin
		_, r.err = r.buf.ReadFrom(os.Stdin)
		r.hasRead = true
	default:
		// read file
		if len(src) > 3 {
			if strings.HasPrefix(src, "@") {
				r.readfile(src[1:])
			} else if r.DefaultAsFile && fsutil.IsFile(src) {
				r.readfile(src)
			}
		}

		if r.err == nil && r.typ == "" {
			r.directToBuf()
		}
	}
}

// HasRead bool
func (r *SourceReader) directToBuf() {
	r.typ = TypeString
	r.buf.WriteString(r.src)
}

// HasRead bool
func (r *SourceReader) readfile(fpath string) {
	fh, err := fsutil.OpenReadFile(fpath)
	if err != nil {
		r.err = err
		return
	}

	r.typ = TypeFile
	_, r.err = r.buf.ReadFrom(fh)
	r.hasRead = true
	fh.Close()
}

// HasRead bool
func (r *SourceReader) HasRead() bool {
	return r.hasRead
}

// SrcType get
func (r *SourceReader) SrcType() string {
	return r.typ
}

// Type get
func (r *SourceReader) Type() string {
	return r.typ
}

// Err get
func (r *SourceReader) Err() error {
	return r.err
}
