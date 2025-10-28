//go:build windows
package util

import (
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestNormalizePath(t *testing.T) {
	testutil.MockEnvValues(testutil.M{"USERPROFILE": "/home/user1", "SHELL": "/bin/bash"}, func() {
		// expand home directory
		assert.Eq(t, "/home/user1/bin", NormalizePath("~/bin"))
		// windows bash 特殊处理
		assert.True(t, IsHookBash())
		assert.Eq(t, "/d/tools/bin", NormalizePath("D:\\tools\\bin"))
	})
}

