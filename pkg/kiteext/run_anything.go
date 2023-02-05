package kiteext

import "github.com/gookit/goutil/sysutil/cmdr"

// RunAnything struct
type RunAnything struct {
}

func NewRunAnything() *RunAnything {
	return &RunAnything{}
}

func (r *RunAnything) Run(name string, args []string) error {
	return nil
}

func (r *RunAnything) RunScript(name string, args []string) error {
	return nil
}

func (r *RunAnything) RunKiteCmd(name string, args []string) error {
	return nil
}

func (r *RunAnything) RunSyscmd(name string, args []string) error {
	return cmdr.NewCmd(name, args...).FlushRun()
}
