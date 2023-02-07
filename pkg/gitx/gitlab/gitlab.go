package gitlab

import (
	"github.com/gookit/goutil/maputil"
	"github.com/inhere/kite/pkg/gitx"
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
	// *gitx.Config // TODO pkg mapstructure not support set value

	// GitUrl git host url
	GitUrl string `json:"git_url"`
	// HostUrl http host url
	HostUrl string `json:"host_url"`
	// AutoSign auto add author sign on commit
	AutoSign bool `json:"auto_sign"`
	// ForkMode enable git fork mode for develop.
	// If is False, use branch mode, will ignore SourceRemote setting.
	ForkMode bool `json:"fork_mode"`
	// SourceRemote the source remote name, it is center repo.
	SourceRemote string `json:"source_remote"`
	// DefaultRemote the default upstream remote name, use for develop. default: origin.
	// It should be forked from SourceRemote.
	DefaultRemote string `json:"default_remote"`
	// BranchAliases branch aliases
	BranchAliases maputil.Aliases `json:"branch_aliases"`

	// ApiUrl api url
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
func New() *GitLab {
	return &GitLab{
		SourceRemote:  gitx.DefaultSrcRemote,
		DefaultRemote: gitx.DefaultOriRemote,
	}
}

// ResolveAlias branch name
func (g *GitLab) ResolveAlias(name string) string {
	return g.BranchAliases.ResolveAlias(name)
}

// LocGlProject instance
func (g *GitLab) LocGlProject(dir string) *GlProject {
	return NewGlProject(dir, g)
}
