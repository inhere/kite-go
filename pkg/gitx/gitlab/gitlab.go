package gitlab

import (
	"github.com/gookit/gitw"
	"github.com/gookit/goutil/maputil"
)

const (
	ProjIdSep  = "%2F"
	DefMainRmt = "main"
)

// GitLab struct
//
// Gen by:
//
//	kite go gen st -s @c -t yml --name GitLab
type GitLab struct {
	// HostUrl host url
	HostUrl string `json:"host_url"`
	// GitUrl git url
	GitUrl string `json:"git_url"`
	// ApiUrl api url
	BaseApi string `json:"base_api"`
	// Token person token.
	// from /profile/personal_access_tokens
	Token string `json:"token"`
	// UpstreamRemote the main remote address name
	UpstreamRemote string `json:"main_remote"`
	// DefaultRemote fork remote name, it's default remote, use for develop.
	DefaultRemote string `json:"fork_remote"`
	// BranchAliases branch aliases
	BranchAliases maputil.Aliases `json:"branch_aliases"`
	// DenyBranches deny as source branch for create PR.
	DenyBranches map[string]string `json:"deny_branches"`
}

// New instance.
func New() *GitLab {
	return &GitLab{
		UpstreamRemote: DefMainRmt,
		DefaultRemote:  gitw.DefaultRemoteName,
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
