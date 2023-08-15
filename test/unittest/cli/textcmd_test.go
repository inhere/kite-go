package cli_test

import (
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

func TestCmd_text_replace(t *testing.T) {
	txtFile := tdataDir + "/text/replace.txt"
	_, err := fsutil.PutContents(txtFile, "hello world")
	assert.NoError(t, err)

	st := app.Cli().RunLine("txt replace -f hello -t hi " + txtFile)
	assert.Eq(t, st, 0)
}
