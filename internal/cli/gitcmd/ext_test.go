package gitcmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/gitx"
)

func TestOpenRemoteReturnsErrorOnMissingRemote(t *testing.T) {
	workdir := initGitRepoForOpenRemote(t)
	oldOpts := orOpts
	orOpts = struct {
		source   bool
		remote   string
		repoPath string
	}{}
	t.Cleanup(func() {
		orOpts = oldOpts
	})

	cmd := NewOpenRemoteCmd(func() *gitx.Config {
		return &gitx.Config{
			HostUrl:      "https://gitlab.example.com",
			SourceRemote: "main",
		}
	})
	initOpenRemoteCmd(t, cmd, workdir)

	err := cmd.Run([]string{"--source"})
	if !assert.Err(t, err) {
		return
	}
	assert.Contains(t, err.Error(), `remote "main" is not found`)
}

func TestOpenRemoteWithRepoPathReturnsErrorOnMissingRemote(t *testing.T) {
	workdir := initGitRepoForOpenRemote(t)
	oldOpts := orOpts
	orOpts = struct {
		source   bool
		remote   string
		repoPath string
	}{}
	t.Cleanup(func() {
		orOpts = oldOpts
	})

	cmd := NewOpenRemoteCmd(func() *gitx.Config {
		return &gitx.Config{}
	})
	initOpenRemoteCmd(t, cmd, workdir)

	err := cmd.Run([]string{"--remote", "main", "group/repo"})
	if !assert.Err(t, err) {
		return
	}
	assert.Contains(t, err.Error(), `remote "main" is not found`)
}

func TestOpenRemoteReturnsParseErrorOnInvalidRemoteURL(t *testing.T) {
	workdir := initGitRepoForOpenRemoteWithURL(t, "invalid-url")
	oldOpts := orOpts
	orOpts = struct {
		source   bool
		remote   string
		repoPath string
	}{}
	t.Cleanup(func() {
		orOpts = oldOpts
	})

	cmd := NewOpenRemoteCmd(func() *gitx.Config {
		return &gitx.Config{}
	})
	initOpenRemoteCmd(t, cmd, workdir)

	err := cmd.Run([]string{"group/repo"})
	if !assert.Err(t, err) {
		return
	}
	assert.Contains(t, err.Error(), `remote "origin" is invalid`)
	assert.Contains(t, err.Error(), "invalid http URL path")
}

func initGitRepoForOpenRemote(t *testing.T) string {
	return initGitRepoForOpenRemoteWithURL(t, "https://gitlab.example.com/group/repo.git")
}

func initGitRepoForOpenRemoteWithURL(t *testing.T, remoteURL string) string {
	t.Helper()

	dir := t.TempDir()
	runGitForOpenRemote(t, dir, "init")
	runGitForOpenRemote(t, dir, "remote", "add", "origin", remoteURL)
	return dir
}

func initOpenRemoteCmd(t *testing.T, cmd interface {
	Init()
	ChWorkDir(string) error
}, dir string) {
	t.Helper()

	oldWd, err := os.Getwd()
	assert.NoErr(t, err)
	cmd.Init()
	assert.NoErr(t, cmd.ChWorkDir(dir))
	t.Cleanup(func() {
		assert.NoErr(t, os.Chdir(oldWd))
	})
}

func runGitForOpenRemote(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
