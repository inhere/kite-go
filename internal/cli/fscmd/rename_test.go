package fscmd

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func Test_formatRenamePattern(t *testing.T) {
	assert.Eq(t, `^(\w+)\.go`, formatRenamePattern("^{word}\\.go"))
	assert.Eq(t, `^(\w+)-(\w+)`, formatRenamePattern("^{word}-{word2}"))
}

func TestHandleRename(t *testing.T) {
	renameOpts := &RenameOptions{
		// Dirs:        ".",
		DryRun:      true,
		Verbose:     true,
		Pattern:     "^(.*)\\.go",
		Replacement: "$1.new.go",
	}

	renameOpts.paths = []string{"*.go"}
	err := handleRename(renameOpts)
	assert.NoError(t, err)
}
