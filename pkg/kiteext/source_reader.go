package kiteext

// SourceReader struct
type SourceReader struct {
	src string
}

// NewSourceReader instance
func NewSourceReader(src string) *SourceReader {
	return &SourceReader{
		src: src,
	}
}

func (r *SourceReader) ReadString() string {
	return r.src
}
