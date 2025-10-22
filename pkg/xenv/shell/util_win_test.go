//go:build windows
package shell

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestNormalizePath(t *testing.T) {
	testutil.MockEnvValues(testutil.M{"USERPROFILE": "/home/user1", "SHELL": "/bin/bash"}, func() {
		// expand home directory
		assert.Eq(t, "/home/user1/bin", NormalizePath("~/bin"))
		// windows bash 特殊处理
		assert.True(t, IsHookWinBash())
		assert.Eq(t, "/d/tools/bin", NormalizePath("D:\\tools\\bin"))
	})
}

func TestProfilePath(t *testing.T) {
	ret, err := sysutil.ExecCmd("echo", []string{"$PROFILE.CurrentUserAllHosts"})
	assert.NoErr(t, err)
	fmt.Println(ret)
}
