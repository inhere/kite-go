package gitx_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/gitx"
)

func TestGlobMatch_Match(t *testing.T) {
	m := gitx.NewBranchMatcher("fea*", false)
	assert.True(t, m.Match("fea-1"))
	assert.True(t, m.Match("fea_dev"))
	assert.False(t, m.Match("fix_2"))

	m = gitx.NewBranchMatcher("fix*", false)
	assert.False(t, m.Match("fea-1"))
	assert.False(t, m.Match("fea_dev"))
	assert.True(t, m.Match("fix_2"))

	m = gitx.NewBranchMatcher("ma*", false)
	assert.True(t, m.Match("main"))
	assert.True(t, m.Match("master"))
	assert.False(t, m.Match("x-main"))
}
