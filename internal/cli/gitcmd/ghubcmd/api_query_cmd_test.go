package ghubcmd

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestApiLatestCommitOptions_requireRepo(t *testing.T) {
	opts := apiLatestCommitOptions{}

	_, err := opts.toInput()
	assert.Err(t, err)
	assert.Contains(t, err.Error(), "repository path")
}

func TestApiTagListOptions_defaultLimit(t *testing.T) {
	opts := apiTagListOptions{
		RepoPath: "owner/repo",
	}

	in, err := opts.toInput()
	assert.NoErr(t, err)
	assert.Eq(t, 20, in.Limit)
}

func TestApiTagListOptions_capLimit(t *testing.T) {
	opts := apiTagListOptions{
		RepoPath: "owner/repo",
		Limit:    200,
	}

	in, err := opts.toInput()
	assert.NoErr(t, err)
	assert.Eq(t, 100, in.Limit)
}
