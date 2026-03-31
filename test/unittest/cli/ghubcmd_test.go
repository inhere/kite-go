package cli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/internal/app"
	"github.com/inhere/kite-go/internal/cli/gitcmd/ghubcmd"
)

func TestCmd_gh_api_tag_add_registered(t *testing.T) {
	cmd := app.Cli.MatchByPath("github:api:tag:add")
	assert.NotNil(t, cmd)
	assert.Eq(t, "add", cmd.Name)
}

func TestCmd_gh_api_commit_latest_registered(t *testing.T) {
	cmd := app.Cli.MatchByPath("github:api:commit:latest")
	assert.NotNil(t, cmd)
	assert.Eq(t, "latest", cmd.Name)
}

func TestCmd_gh_api_tag_list_registered(t *testing.T) {
	cmd := app.Cli.MatchByPath("github:api:tag:list")
	assert.NotNil(t, cmd)
	assert.Eq(t, "list", cmd.Name)
}

func TestCmd_gh_api_tag_add_require_repo(t *testing.T) {
	opts := &ghubcmd.ApiTagAddOptionsForTest
	*opts = ghubcmd.ApiTagAddOptionsForTestType{
		Version: "v1.2.3",
		Message: "test-message",
		SHA:     "abcdef1234567890",
	}

	_, err := opts.ToCreateInputForTest()
	assert.Err(t, err)
	assert.Contains(t, err.Error(), "repository path")
}

func TestCmd_gh_api_tag_add_allow_empty_sha(t *testing.T) {
	opts := &ghubcmd.ApiTagAddOptionsForTest
	*opts = ghubcmd.ApiTagAddOptionsForTestType{
		RepoPath: "owner/repo",
		Version:  "v1.2.3",
		Message:  "test-message",
	}

	in, err := opts.ToCreateInputForTest()
	assert.NoErr(t, err)
	assert.Eq(t, "", in.Object)
}

func TestCmd_gh_api_tag_add_normalize_version(t *testing.T) {
	opts := ghubcmd.ApiTagAddOptionsForTestType{
		RepoPath: "owner/repo",
		Version:  "1.2.3",
		Message:  "test-message",
		SHA:      "abcdef1234567890",
	}

	in, err := opts.ToCreateInputForTest()
	assert.NoErr(t, err)
	assert.Eq(t, "v1.2.3", in.Version)
}

var _ *gcli.Command
