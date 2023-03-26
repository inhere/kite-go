package quickjump_test

import (
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/quickjump"
)

func TestQuickJump_Save(t *testing.T) {
	qj := quickjump.NewQuickJump()
	qj.DataDir = "./testdata"

	err := qj.Init()
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(qj.Datafile()))

	qj.AddNamed("name1", "path1")
	qj.AddNamedPaths(map[string]string{
		"name2": "path2/sub2",
		"name3": "/path3/to/sub3",
		"home":  "/path3/to/home",
	})
	assert.Eq(t, "/path3/to/sub3", qj.Match("name3"))

	ss := qj.Search([]string{"path3", "home"}, 3)
	assert.NotEmpty(t, ss)
	assert.Len(t, ss, 1)
	assert.Eq(t, "/path3/to/home", ss[0])

	qj.AddHistory("path4/to/sub4")
	assert.NotEmpty(t, qj.Histories)
	qj.AddHistory("/path5/to/sub5")
}
