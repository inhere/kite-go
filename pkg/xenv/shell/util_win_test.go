//go:build windows
package shell

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestProfilePath(t *testing.T) {
	ret, err := sysutil.ExecCmd("echo", []string{"$PROFILE.CurrentUserAllHosts"})
	assert.NoErr(t, err)
	fmt.Println(ret)
}
