package cli_test

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
)

// test for textcmd.NewTemplateCmd(true)
func TestCmd_fs_render(t *testing.T) {
	txtFile := tdataDir + "/fs/render.txt"
	tplText := "hi, my name is {{ name }}, age is {{age}}"
	fsutil.MustSave(txtFile, tplText)

	st := app.Cli().RunLine("fs render -v name=Tom -v age=18 -w " + txtFile)
	assert.Eq(t, st, 0)

	s := fsutil.ReadString(txtFile)
	fmt.Println(s)
	assert.StrContains(t, s, "Tom")
}

// test for fscmd.NewReplaceCmd()
func TestCmd_fs_replace(t *testing.T) {
	txtFile := tdataDir + "/fs/replace.txt"
	_, err := fsutil.PutContents(txtFile, "hello world")
	assert.NoError(t, err)

	st := app.Cli().RunLine("fs replace -f hello -t hi " + txtFile)
	assert.Eq(t, st, 0)
}
