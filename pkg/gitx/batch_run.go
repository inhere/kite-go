package gitx

func BatchPull() {

}

// GitBatchRun struct
type GitBatchRun struct {
}

func NewBatchRun(fn ...func(gbr *GitBatchRun)) *GitBatchRun {
	gbr := &GitBatchRun{}

	if len(fn) > 0 {
		fn[0](gbr)
	}
	return gbr
}
