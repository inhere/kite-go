package gituse

import "github.com/inherelab/kite/pkg/cmdutil"

func BatchPull() {

}

// GitBatchRun
type GitBatchRun struct {
	cmdutil.CmdRunner
}

func NewBatchRun(fn ...func(gbr *GitBatchRun)) *GitBatchRun {
	gbr := &GitBatchRun{}

	if len(fn) > 0 {
		fn[0](gbr)
	}
	return gbr
}

