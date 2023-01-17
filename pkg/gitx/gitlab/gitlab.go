package gitlab

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
	// MainRemote main remote
	MainRemote string `json:"main_remote"`
	// ForkRemote fork remote
	ForkRemote string `json:"fork_remote"`
	// BranchAliases branch aliases
	BranchAliases map[string]string `json:"branch_aliases"`
}

func (g *GitLab) DefaultRemote() string {
	return g.ForkRemote
}
