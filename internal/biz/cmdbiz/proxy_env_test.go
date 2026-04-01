package cmdbiz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyEnv(t *testing.T) {
	pcc := &ProxyCmdConf{
		CommandIds: []string{"github:tag:*"},
	}

	isMatch := pcc.IsMatchName("github", "tag", []string{"del", "v0.2.0"})
	assert.False(t, isMatch)

	pcc.CommandIds = []string{"github:tag:**"}
	isMatch = pcc.IsMatchName("github", "tag", []string{"del", "v0.2.0"})
	assert.True(t, isMatch)
}
