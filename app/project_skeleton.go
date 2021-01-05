package app

// Skeleton project skeleton struct
type Skeleton struct {
	file string
	Name string
	Desc string
	//
	PkgName string
	Defines map[string]interface{}
}

func NewSkeleton(file string) *Skeleton {
	s := &Skeleton{file: file}

	return s
}

