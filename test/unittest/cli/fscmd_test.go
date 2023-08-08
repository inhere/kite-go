package cli_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

func TestCmd_fs_render(t *testing.T) {
	st := app.Cli().RunLine("fs render -v name=Tom -v age=18 ../../testdata/fs/render-01.tpl")
	assert.Eq(t, st, 0)
}
