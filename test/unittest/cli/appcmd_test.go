package cli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

// see appcmd.NewPathMapCmd()
func TestCmd_app_pathMap(t *testing.T) {
	gcli.SetVerbose(gcli.VerbWarn)

	st := app.Cli().RunLine("app pmp @pwd")
	assert.Eq(t, st, 0)

	st = app.Cli().RunLine("app pmp @home")
	assert.Eq(t, st, 0)
}
