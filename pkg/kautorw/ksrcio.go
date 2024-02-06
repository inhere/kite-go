package kautorw

// types for read or write
const (
	DstStdin  = "@stdin"
	DstStdout = "@stdout"
	DstClip   = "@clip"
	DstSrc    = "@src" // src file path
)

// types for read or write
const (
	TypeStdin  = "stdin"
	TypeString = "string" // raw string
	TypeStdout = "stdout"

	TypeClip = "clip"
	TypeFile = "file"
)

func ReadStdin(fns ...ReaderFn) (string, error) {
	return NewSourceReader(DstStdin, fns...).ReadStdin().TryString()
}

func ReadClip(fns ...ReaderFn) (string, error) {
	return NewSourceReader(DstClip, fns...).ReadClip().TryString()
}

func ReadContents(src string, fns ...ReaderFn) (string, error) {
	return NewSourceReader(src, fns...).TryReadString()
}

func WriteContents(contents, dst string) (string, error) {
	return "", nil
}
