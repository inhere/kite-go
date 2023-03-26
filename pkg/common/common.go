package common

// PathResolver struct
type PathResolver struct {
	// PathResolve handler
	PathResolve func(path string) string `json:"-"`
}
