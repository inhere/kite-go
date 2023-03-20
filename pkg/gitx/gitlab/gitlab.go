package gitlab

import (
	"github.com/inhere/kite-go/pkg/gitx"
)

const (
	ProjectIdSep = "%2F"
)

// GitLab config struct for gitlab
//
// Gen by:
//
//	kite go gen st -s @c -t yml --name GitLab
type GitLab struct {
	*gitx.Config

	// BaseApi api url
	BaseApi string `json:"base_api"`
	// Token person token.
	// - from /profile/personal_access_tokens
	Token string `json:"token"`
	// BranchAliases branch aliases
	// BranchAliases maputil.Aliases `json:"branch_aliases"`
	// DenyBranches deny as source branch for create PR.
	DenyBranches map[string]string `json:"deny_branches"`
}

// New instance.
func New(cfg *gitx.Config) *GitLab {
	cfg.HostType = gitx.HostGitlab

	return &GitLab{
		Config: cfg,
	}
}

// LocGlProject instance
func (g *GitLab) LocGlProject(dir string) *GlProject {
	return NewGlProject(dir, g)
}
