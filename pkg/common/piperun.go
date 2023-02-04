package common

import (
	"github.com/gookit/goutil/structs"
)

// RunFn func
type RunFn func(ctx *structs.Data) error

// PipeRun struct
type PipeRun struct {
	ctx *structs.Data
	err error
	fns []RunFn
}

func NewPipeRun() *PipeRun {
	return &PipeRun{
		ctx: structs.NewData(),
	}
}

func (p *PipeRun) Add(fn RunFn) *PipeRun {

	return p
}

func (p *PipeRun) Run() error {
	for _, fn := range p.fns {
		if err := fn(p.ctx); err != nil {
			return err
		}
	}
	return nil
}
