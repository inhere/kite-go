package kiteext

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil/clipboard"
)

// types for read or write
const (
	DstStdin  = "@stdin"
	DstStdout = "@stdout"
	DstClip   = "@clip"
)

// types for read or write
const (
	TypeStdin  = "stdin"
	TypeString = "string" // raw string

	TypeStdout = "stdout"

	TypeClip = "clip"
	TypeFile = "file"
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
	// emptyAction type on src is empty.
	emptyAction string
	// has read contents by src
	hasRead   bool
	trimSpace bool
	// return error on result is empty
	checkResult bool
}

func ReadStdin(fns ...ReaderFn) (string, error) {
	return NewSourceReader(DstStdin, fns...).ReadStdin().TryString()
}

func ReadClip(fns ...ReaderFn) (string, error) {
	return NewSourceReader(DstClip, fns...).ReadClip().TryString()
}

func ReadContents(src string, fns ...ReaderFn) (string, error) {
	return NewSourceReader(src, fns...).TryReadString()
}

// WithEmptyAction setting.
func WithEmptyAction(typ string) ReaderFn {
	return func(sr *SourceReader) { sr.emptyAction = typ }
}

// FallbackStdin read from stdin.
func FallbackStdin() ReaderFn {
	return func(sr *SourceReader) { sr.emptyAction = TypeStdin }
}

// WithTrimSpace on read contents.
func WithTrimSpace() ReaderFn {
	return func(sr *SourceReader) { sr.trimSpace = true }
}

// WithCheckResult on read contents.
func WithCheckResult() ReaderFn {
	return func(sr *SourceReader) { sr.checkResult = true }
}

// NewSourceReader instance
func NewSourceReader(src string, fns ...ReaderFn) *SourceReader {
	sr := &SourceReader{
		src: src,
		typ: TypeString,
	}

	for _, fn := range fns {
		fn(sr)
	}
	return sr
}

// SetEmptyAction type on src is empty
func (r *SourceReader) SetEmptyAction(typ string) *SourceReader {
	r.emptyAction = typ
	return r
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
		r.buf.WriteString(r.src)
		return
	}

	switch r.src {
	case "@c", "@cb", "@clip", "@clipboard", "clipboard":
		r.typ = TypeClip
		r.err = clipboard.Std().ReadTo(&r.buf)
		r.hasRead = true
	default:
		r.buf.WriteString(r.src)
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
		r.buf.WriteString(r.src)
		return
	}

	switch r.src {
	case "@in", "@stdin", "stdin":
		r.typ = TypeStdin
		_, r.err = r.buf.ReadFrom(os.Stdin)
		r.hasRead = true
	default:
		r.buf.WriteString(r.src)
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
	if r.trimSpace {
		s = strings.TrimSpace(s)
	}

	if r.checkResult && len(s) == 0 {
		r.err = emptyResultErr
	}
	return s
}

// tryReadAny return string
func (r *SourceReader) tryReadAny() {
	ln := len(r.src)
	if ln > 86 {
		r.buf.WriteString(r.src)
		return
	}

	src := r.src
	if r.emptyAction != "" && ln == 0 {
		src = "@" + r.emptyAction
	}

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
		if ln > 3 && strings.HasPrefix(r.src, "@") {
			r.readfile(r.src[1:])
		} else {
			r.buf.WriteString(r.src)
		}
	}
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
