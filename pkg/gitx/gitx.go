package gitx

import (
	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
)

// GitLoc repo struct
type GitLoc struct {
	*Config
	// local repo
	repo *gitw.Repo
}

// NewGitLoc instance
func NewGitLoc(repoDir string, cfg *Config) *GitLoc {
	if cfg == nil {
		cfg = NewConfig()
	} else {
		cfg.Init()
	}

	return &GitLoc{
		repo: gitw.NewRepo(repoDir).PrintCmdOnExec(),
		// config
		Config: cfg,
	}
}

func (g *GitLoc) HasDefaultRemote() bool {
	return g.Repo().HasRemote(g.DefaultRemote)
}

func (g *GitLoc) HasSourceRemote() bool {
	return g.Repo().HasRemote(g.SourceRemote)
}

// DefRemoteInfo data.
func (g *GitLoc) DefRemoteInfo() *gitw.RemoteInfo {
	ri := g.Repo().RemoteInfo(g.DefaultRemote)
	if ri != nil {
		goutil.Panicf("gitx: default remote %q is not found", g.DefaultRemote)
	}
	return ri
}

// SrcRemoteInfo data.
func (g *GitLoc) SrcRemoteInfo() *gitw.RemoteInfo {
	ri := g.Repo().RemoteInfo(g.SourceRemote)
	if ri != nil {
		goutil.Panicf("gitx: main repo remote %q is not found", g.SourceRemote)
	}
	return ri
}

// Check git config.
func (g *GitLoc) Check() error {
	return nil
}

// Repo instance.
func (g *GitLoc) Repo() *gitw.Repo {
	return g.repo
}

// RepoDir path.
func (g *GitLoc) RepoDir() string {
	return g.repo.Dir()
}

// Cmd create
func (g *GitLoc) Cmd(name string, args ...string) *gitw.GitWrap {
	return g.repo.Cmd(name, args...)
}
