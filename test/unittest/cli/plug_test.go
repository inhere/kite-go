package cli

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

// test for extcmd.PlugCmd
func TestCmd_plug_run01(t *testing.T) {
	gcli.SetVerbose(gcli.VerbWarn)

	st := app.Cli.RunLine("plug @base/plugins/test-plug.go --name inhere --age 18")
	assert.Eq(t, st, 0)
}
