package cli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

// see xenvcmd.XEnvCmd
func TestCmd_xenv_list_activity(t *testing.T) {
	gcli.SetVerbose(gcli.VerbWarn)

	st := app.Cli.RunLine("xenv -d list act")
	assert.Eq(t, st, 0)

}

func TestCmd_xenv_list_sdks(t *testing.T) {
	gcli.SetVerbose(gcli.VerbWarn)

	t.Run("in bash", func(t *testing.T) {
		testutil.MockEnvValues(map[string]string{
			"XENV_HOOK_SHELL": "bash",
		}, func() {
			st := app.Cli.RunLine("xenv list tools")
			assert.Eq(t, st, 0)
		})
	})

	t.Run("in pwsh", func(t *testing.T) {
		testutil.MockEnvValues(map[string]string{
			"XENV_HOOK_SHELL": "pwsh",
		}, func() {
			st := app.Cli.RunLine("xenv list tools")
			assert.Eq(t, st, 0)
		})
	})
}
