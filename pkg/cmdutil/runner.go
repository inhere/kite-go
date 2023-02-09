package cmdutil

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/goutil/sysutil/cmdr"
)

// Runner struct
type Runner struct {
	*cmdr.Runner

	// Silent not print messages
	Silent  bool
	Confirm bool
}

// NewRunner instance
func NewRunner(fns ...func(rr *Runner)) *Runner {
	rr := &Runner{
		Runner: cmdr.NewRunner(),
	}

	for _, fn := range fns {
		fn(rr)
	}
	return rr
}

// Run all tasks
func (r *Runner) Run() error {
	if !r.Silent {
		color.Magentaf("\n# Run All Tasks(%d steps):\n", r.Len())
	}

	r.BeforeRun = func(cr *cmdr.Runner, t *cmdr.Task) bool {
		if r.Confirm {
			if !interact.Unconfirmed("continue run?", true) {
				return false
			}
		}

		if !r.Silent {
			color.Yellowf("Step#%d> %s\n", t.Index()+1, t.Cmdline())
		}
		return true
	}

	err := r.Runner.Run()

	if !r.Silent && err == nil {
		color.Successln("✅  All Tasks Run Completed")
	}
	return err
}
