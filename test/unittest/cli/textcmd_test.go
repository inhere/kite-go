package cli_test

import (
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

func TestCmd_text_replace(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		txtFile := tdataDir + "/text/replace.txt"
		_, err := fsutil.PutContents(txtFile, "hello world")
		assert.NoError(t, err)

		st := app.Cli.RunLine("txt replace -f hello -t hi @" + txtFile)
		assert.Eq(t, st, 0)
	})

	t.Run("char01", func(t *testing.T) {
		testutil.RewriteStdout()
		st := app.Cli.RunLine("txt repl -f / -t . ab/cd")
		assert.Eq(t, st, 0)
		assert.Eq(t, "ab.cd", testutil.RestoreStdout())
	})
	t.Run("char02", func(t *testing.T) {
		testutil.RewriteStdout()
		st := app.Cli.RunLine("txt repl -f - -t . ab-cd")
		assert.Eq(t, st, 0)
		assert.Eq(t, "ab.cd", testutil.RestoreStdout())
	})
	t.Run("char03", func(t *testing.T) {
		testutil.RewriteStdout()
		st := app.Cli.RunLine("txt repl -f - -t / ab-cd")
		assert.Eq(t, st, 0)
		assert.Eq(t, "ab/cd", testutil.RestoreStdout())
	})

	t.Run("sep_SLASH", func(t *testing.T) {
		testutil.RewriteStdout()
		st := app.Cli.RunLine("txt repl -f - -t SLASH ab-cd")
		assert.Eq(t, st, 0)
		assert.Eq(t, "ab/cd", testutil.RestoreStdout())
	})

}
