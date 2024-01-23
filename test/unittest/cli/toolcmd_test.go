package cli

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

// test for syscmd.NewBatchRunCmd()
func TestCmd_tool_brun(t *testing.T) {
	gcli.SetVerbose(gcli.VerbWarn)

	st := app.Cli().RunLine("tool brun -c 'echo {item}' --for testing,qa,pre")
	assert.Eq(t, st, 0)
}
