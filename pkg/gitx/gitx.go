package gitx

import (
	"github.com/gookit/gitw"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
)

// GitLoc repo struct
type GitLoc struct {
	*Config
	// local repo
	*gitw.Repo
}

// NewGitLoc instance
func NewGitLoc(repoDir string, cfg *Config) *GitLoc {
	if cfg == nil {
		cfg = NewConfig()
	} else {
		cfg.Init()
	}

	return &GitLoc{
		Repo: gitw.NewRepo(repoDir).PrintCmdOnExec(),
		// config
		Config: cfg,
	}
}

// FetchOrigin fetch default remote
func (g *GitLoc) FetchOrigin() error {
	return g.Cmd("fetch", g.DefaultRemote, "-np").Run()
}

// FetchSource fetch source remote
func (g *GitLoc) FetchSource() error {
	return g.Cmd("fetch", g.SourceRemote, "-np").Run()
}

// HasDefaultBranch check
func (g *GitLoc) HasDefaultBranch(br string) bool {
	return g.HasOriginBranch(br)
}

// HasOriginBranch check
func (g *GitLoc) HasOriginBranch(br string) bool {
	return g.HasRemoteBranch(br, g.DefaultRemote)
}

// HasSourceBranch check
func (g *GitLoc) HasSourceBranch(br string) bool {
	return g.HasRemoteBranch(br, g.SourceRemote)
}

func (g *GitLoc) HasDefaultRemote() bool {
	return g.Repo.HasRemote(g.DefaultRemote)
}

func (g *GitLoc) HasSourceRemote() bool {
	return g.Repo.HasRemote(g.SourceRemote)
}

// DefRemoteInfo data.
func (g *GitLoc) DefRemoteInfo() *gitw.RemoteInfo {
	ri := g.Repo.RemoteInfo(g.DefaultRemote)
	if ri != nil {
		goutil.Panicf("gitx: default remote %q is not found", g.DefaultRemote)
	}
	return ri
}

// SrcRemoteInfo data.
func (g *GitLoc) SrcRemoteInfo() *gitw.RemoteInfo {
	ri := g.Repo.RemoteInfo(g.SourceRemote)
	if ri != nil {
		goutil.Panicf("gitx: main repo remote %q is not found", g.SourceRemote)
	}
	return ri
}

// Check git config.
func (g *GitLoc) Check() error {
	if err := g.CheckRemote(); err != nil {
		return err
	}

	// TODO check others
	return nil
}

// CheckRemote git config.
func (g *GitLoc) CheckRemote() error {
	if !g.HasRemote(g.DefaultRemote) {
		return errorx.Rawf("the default remote %q is not exists", g.DefaultRemote)
	}

	if g.ForkMode && !g.HasRemote(g.SourceRemote) {
		return errorx.Rawf("the source remote %q is not exists", g.SourceRemote)
	}

	return nil
}

// RepoDir path.
func (g *GitLoc) RepoDir() string {
	return g.Repo.Dir()
}
