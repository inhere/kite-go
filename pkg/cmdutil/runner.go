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
		color.Magentaf("# Run All Tasks(%d steps):\n", r.Len())
	}

	err := r.Runner.Run()

	if !r.Silent && err == nil {
		color.Successln("Run Completed")
	}

	return err
}

// RunTask command
func (r *Runner) RunTask(task *cmdr.Task) bool {
	if r.Confirm {
		if interact.Unconfirmed("continue run?", true) {
			return true
		}
	}

	if !r.Silent {
		color.Infof("Task #%d: %s", task.Index()+1, task.Cmdline())
	}

	return r.Runner.RunTask(task)
}
